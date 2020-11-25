package upload

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/dathan/go-web-backend/pkg/entities"
	localresponse "github.com/dathan/go-web-backend/pkg/http/response"
	"github.com/davecgh/go-spew/spew"
)

type csvUp struct {
}

func CSVUpload(response *goyave.Response, request *goyave.Request) {
	// authentication happens at a layer above this, this we can assume the logged in user is available
	user := request.User.(*entities.User)

	reader, err := request.Request().MultipartReader()
	if err != nil {
		resp := localresponse.NewResponse(false)
		resp.ErrorMessage = fmt.Sprintf("validationError: %s", err.Error())
		response.JSON(http.StatusNotAcceptable, resp)
		return
	}

	csvUp := &csvUp{}

	cPaths, err := csvUp.saveUploadToDatabase(reader, user)
	if err != nil {
		resp := localresponse.NewResponse(false)
		resp.ErrorMessage = fmt.Sprintf("validationError: %s", err.Error())
		response.JSON(http.StatusNotAcceptable, resp)
		return

	}

	//todo: refactopr below
	for _, row := range cPaths {
		cParsedRow, err := csvUp.validateCSV(*row, user)
		if err != nil {
			resp := localresponse.NewResponse(false)
			resp.ErrorMessage = fmt.Sprintf("validationError: %s", err.Error())
			response.JSON(http.StatusNotAcceptable, resp)
			return
		}

		for _, cPRow := range cParsedRow {
			tx := database.GetConnection().Create(&cPRow)
			if tx.Error != nil {
				resp := localresponse.NewResponse(false)
				resp.ErrorMessage = fmt.Sprintf("DBError: %s", tx.Error.Error())
				response.JSON(http.StatusNotAcceptable, resp)
				return
			}
		}
	}

	resp := localresponse.NewResponse(true)
	response.JSON(http.StatusOK, resp)

}

func CSVList(response *goyave.Response, request *goyave.Request) {
	// authentication happens at a layer above this, this we can assume the logged in user is available
	user := request.User.(*entities.User)
	rows := []entities.Contacts_Paths{}
	tx := database.GetConnection().Where("owner_id = ?", user.ID).Find(&rows)
	if tx.Error != nil {
		resp := localresponse.NewResponse(false)
		resp.ErrorMessage = tx.Error.Error()
		response.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := localresponse.NewResponse(true)
	resp.Contacts_Paths = &rows
}

//
// saveUploadToDatabase takes the files and puts them in the db
func (c *csvUp) saveUploadToDatabase(reader *multipart.Reader, user *entities.User) ([]*entities.Contacts_Paths, error) {
	var files []string
	var cPaths []*entities.Contacts_Paths

	//copy each part to destination.
	for {
		row := &entities.Contacts_Paths{}
		row.OwnerID = user.ID
		part, err := reader.NextPart()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		//if part.FileName() is empty, skip this iteration.
		if part.FileName() == "" {
			continue
		}

		save_dir := "/tmp" // TODO:make better

		fmt.Printf("creating filename: %s\n", save_dir+"/"+part.FileName())

		filename := save_dir + "/" + part.FileName()

		row.FileLocation = filename

		files = append(files, filename)
		dst, err := os.Create(filename)

		if err != nil {
			return nil, err
		}

		defer func() {
			if err := dst.Close(); err != nil {
				fmt.Errorf("Close Error: %s\n", err)
			}

		}()

		byte_data, err := ioutil.ReadAll(part)
		if err != nil {
			return nil, err
		}
		row.CSVData = string(byte_data)
		cPaths = append(cPaths, row)
		if _, err := io.Copy(dst, part); err != nil {
			return nil, err
		}

		tx := database.GetConnection().Create(row)
		if tx.Error != nil {
			return nil, tx.Error
		}

	}

	fmt.Printf("Uploaded %d files\n\t%+v\n", len(files), files)

	return cPaths, nil
}

func (c *csvUp) validateCSV(row entities.Contacts_Paths, user *entities.User) ([]entities.Contacts_Parsed, error) {

	r := csv.NewReader(strings.NewReader(row.CSVData))
	r.Comma = ','
	r.Comment = '#'
	r.LazyQuotes = true
	ret := []entities.Contacts_Parsed{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		spew.Dump(record)

		cpCSVRow := &entities.Contacts_Parsed{}
		//email,last,first,address,city, state, zip
		cpCSVRow.Email = record[0]
		cpCSVRow.LastName = record[1]
		cpCSVRow.FirstName = record[2]

		cpCSVRow.OwnerID = user.ID
		//todo: clean up
		if len(record) == 5 {
			cpCSVRow.StreetAddress = record[3]
			cpCSVRow.CityCode = record[4]
			cpCSVRow.ZipCode = record[5]
		}

		if isValidEmail(record[0]) == false {
			return nil, errors.New("INVALID_EMAIL: " + record[0])
		}

		ret = append(ret, *cpCSVRow)

	}
	return ret, nil
}

func isValidEmail(emailAddress string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(emailAddress) < 3 && len(emailAddress) > 254 {
		return false
	}
	return emailRegex.MatchString(emailAddress)
}

func (c *csvUp) getAllParsedFiles(user *entities.User) {

}
