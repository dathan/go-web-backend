package auth

//
// modified version of the jwt_controller.go code
//
import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/config"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/System-Glitch/goyave/v3/helper"
	"github.com/System-Glitch/goyave/v3/lang"
	"github.com/dathan/go-web-backend/pkg/entities"
	userentity "github.com/dathan/go-web-backend/pkg/entities/user"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWTAuthenticator implementation of Authenticator using a JSON Web Token.
type JWTAuthenticator struct{}

var _ auth.Authenticator = (*JWTAuthenticator)(nil) // implements Authenticator

func init() {
	config.Register("auth.jwt.secret", config.Entry{
		Value:            nil,
		Type:             reflect.String,
		IsSlice:          false,
		AuthorizedValues: []interface{}{},
	})
	config.Register("auth.jwt.expiry", config.Entry{
		Value:            300,
		Type:             reflect.Int,
		IsSlice:          false,
		AuthorizedValues: []interface{}{},
	})
	config.Register("auth.jwt.refresh_expiry", config.Entry{
		Value:            3000000,
		Type:             reflect.Int,
		IsSlice:          false,
		AuthorizedValues: []interface{}{},
	})
	config.Register("auth.jwt.refresh_secret", config.Entry{
		Value:            nil,
		Type:             reflect.String,
		IsSlice:          false,
		AuthorizedValues: []interface{}{},
	})

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
	)

}

// Authenticate fetch the user corresponding to the token
// found in the given request and puts the result in the given user pointer.
// If no user can be authenticated, returns false.
//
// The database request is executed based on the model name and the
// struct tag `auth:"username"`.
//
// This implementation is a JWT-based authentication using HMAC SHA256, supporting only one active token.
func (a *JWTAuthenticator) Authenticate(request *goyave.Request, user interface{}) error {

	tokenString, ok := request.BearerToken()
	if tokenString == "" || !ok {
		return fmt.Errorf(lang.Get(request.Lang, "auth.no-credentials-provided"))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return satisfySignedString(token, "auth.jwt.secret")
	})

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			column := auth.FindColumns(user, "username")[0]
			result := database.GetConnection().Where(column.Name+" = ?", claims["username"]).First(user)

			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					return fmt.Errorf(lang.Get(request.Lang, "auth.invalid-credentials"))
				}
				panic(result.Error)
			}

			return nil
		}
	}

	return a.makeError(request.Lang, err.(*jwt.ValidationError).Errors)
}

func (a *JWTAuthenticator) makeError(language string, bitfield uint32) error {
	if bitfield&jwt.ValidationErrorNotValidYet != 0 {
		return fmt.Errorf(lang.Get(language, "auth.jwt-not-valid-yet"))
	} else if bitfield&jwt.ValidationErrorExpired != 0 {
		return fmt.Errorf(lang.Get(language, "auth.jwt-expired"))
	}
	return fmt.Errorf(lang.Get(language, "auth.jwt-invalid"))
}

// login
func (c *JWTAuthenticator) Login(response *goyave.Response, request *goyave.Request) {
	user := userentity.User{}
	username := request.String("username")
	columns := auth.FindColumns(user, "username", "password")
	resp := entities.NewResponse(false)

	result := database.GetConnection().Where(columns[0].Name+" = ?", username).First(&user)
	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !notFound {
		resp.ErrorMessage = result.Error.Error()
		response.JSON(http.StatusNotFound, resp)
		return
	}

	pass := reflect.Indirect(reflect.ValueOf(user)).FieldByName(columns[1].Field.Name)
	if notFound || bcrypt.CompareHashAndPassword([]byte(pass.String()), []byte(request.String("password"))) != nil {
		resp.ErrorMessage = fmt.Sprintf("validationError: %s", lang.Get(request.Lang, "auth.invalid-credentials"))
		response.JSON(http.StatusUnprocessableEntity, resp)
		return
	}

	c.ResponseJWT(&user, resp, response)
}

// Refresh a token is a handler that looks at the refersh token and respond with a new refesh token
func (c *JWTAuthenticator) Refresh(response *goyave.Response, request *goyave.Request) {

	refreshToken := request.String("refresh_token")
	resp := entities.NewResponse(false)

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}

		return satisfySignedString(token, "auth.jwt.refresh_secret")
	})

	if err != nil {
		response.JSON(http.StatusUnauthorized, resp)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		resp.ErrorMessage = err.Error()
		response.JSON(http.StatusUnauthorized, resp)
		return
	}

	// since there token is valid
	// todo: mark the current refresh_token as invalid and generate a new token
	user, err := tokenToUser(token)
	if err != nil {
		resp.ErrorMessage = err.Error()
		response.JSON(http.StatusNotFound, resp)
		return
	}

	c.ResponseJWT(user, resp, response)
}

func (a *JWTAuthenticator) ResponseJWT(user *userentity.User, resp *entities.CommonResponse, response *goyave.Response) {

	tokenStr, err := GenerateToken(user, "auth")
	if err != nil {
		resp.ErrorMessage = err.Error()
		response.JSON(http.StatusNotFound, resp)
		return
	}

	resp.OK = true
	resp.Token = tokenStr

	refresh_token, err := GenerateToken(user, "refresh")
	if err != nil {
		resp.ErrorMessage = err.Error()
		response.JSON(http.StatusNotFound, resp)
		return
	}

	resp.RefreshToken = refresh_token
	resp.User = user
	response.JSON(http.StatusOK, resp)
}

// Google start of the login
func (c *JWTAuthenticator) GoogleLogin(response *goyave.Response, request *goyave.Request) {
	gothic.BeginAuthHandler(response, request)
}

//
func (c *JWTAuthenticator) GoogleAuthCallBack(response *goyave.Response, request *goyave.Request) {
	resp := entities.NewResponse(false)
	user, err := gothic.CompleteUserAuth(response, request)
	if err != nil {
		response.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.OK = true
	spew.Dump(user)
	response.JSON(http.StatusOK, resp)
}

// GenerateToken generate a new JWT.
// The token is created using the HMAC SHA256 method and signed using
// the "auth.jwt.secret" config entry.
// The token is set to expire in the amount of seconds defined by
// the "auth.jwt.expiry" config entry.
func GenerateToken(user *userentity.User, tokenType string) (string, error) {
	var expiry time.Duration
	var expiryKey string = "auth.jwt.expiry"
	var secretKey string = "auth.jwt.secret"

	if tokenType == "refresh" {
		expiryKey = "auth.jwt.refresh_expiry"
		secretKey = "auth.jwt.refresh_secret"
	}

	expiry = time.Duration(config.GetInt(expiryKey)) * time.Second

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.UserName,
		"nbf":      now.Unix(),             // Not Before
		"exp":      now.Add(expiry).Unix(), // Expiry
	})

	return token.SignedString(signedString(user, config.GetString(secretKey)))
}

// satisfySignedString is using my custom strong. token at this stage is not valid. We are using a callback to return the key
func satisfySignedString(token *jwt.Token, config_secret_key string) ([]byte, error) {

	user, err := tokenToUser(token)
	if err != nil {
		return nil, err
	}

	return signedString(user, config.GetString(config_secret_key)), nil

}

func tokenToUser(token *jwt.Token) (*userentity.User, error) {
	var user userentity.User
	var claims jwt.MapClaims
	var ok bool = false

	if claims, ok = token.Claims.(jwt.MapClaims); ok {
		column := auth.FindColumns(user, "username")[0]
		result := database.GetConnection().Where(column.Name+" = ?", claims["username"]).First(&user)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	if user.UserName != claims["username"] {
		return nil, errors.New("claim is not valid") // TODO lang
	}
	return &user, nil

}

func signedString(user *userentity.User, secret string) []byte {
	return []byte(user.Password + ":" + secret)
}

// FindColumns in given struct. A field matches if it has a "auth" tag with the given value.
// Returns a slice of found fields, ordered as the input "fields" slice.
// If the nth field is not found, the nth value of the returned slice will be nil.
//
// Promoted fields are matched as well.
//
// Given the following struct and "username", "notatag", "password":
//  type TestUser struct {
// 		gorm.Model
// 		Name     string `gorm:"type:varchar(100)"`
// 		Password string `gorm:"type:varchar(100)" auth:"password"`
// 		Email    string `gorm:"type:varchar(100);unique_index" auth:"username"`
//  }
//
// The result will be the "Email" field, "nil" and the "Password" field.
func FindColumns(strct interface{}, fields ...string) []*auth.Column {
	length := len(fields)
	result := make([]*auth.Column, length)

	value := reflect.ValueOf(strct)
	t := reflect.TypeOf(strct)
	if t.Kind() == reflect.Ptr {
		value = value.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := value.Field(i)
		fieldType := t.Field(i)
		if field.Kind() == reflect.Struct && fieldType.Anonymous {
			// Check promoted fields recursively
			for i, v := range FindColumns(field.Interface(), fields...) {
				if v != nil {
					result[i] = v
				}
			}
			continue
		}

		tag := fieldType.Tag.Get("auth")
		if index := helper.IndexOf(fields, tag); index != -1 {
			result[index] = &auth.Column{
				Name:  columnName(&fieldType),
				Field: &fieldType,
			}
		}
	}

	return result
}

func columnName(field *reflect.StructField) string {
	for _, t := range strings.Split(field.Tag.Get("gorm"), ";") { // Check for gorm column name override
		if strings.HasPrefix(t, "column") {
			v := strings.Split(t, ":")
			return strings.TrimSpace(v[1])
		}
	}

	return database.Conn().Config.NamingStrategy.ColumnName("", field.Name)
}
