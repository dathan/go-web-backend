package contacts

import (
	"net/http"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/dathan/go-web-backend/pkg/entities"
	localresponse "github.com/dathan/go-web-backend/pkg/http/response"
)

//List provides the full contact list TODO add pagination
func List(response *goyave.Response, request *goyave.Request) {
	// authentication happens at a layer above this, this we can assume the logged in user is available
	user := request.User.(*entities.User)
	rows := []entities.Contacts_Parsed{}
	tx := database.GetConnection().Where("owner_id = ?", user.ID).Find(&rows)
	if tx.Error != nil {
		resp := localresponse.NewResponse(false)
		resp.ErrorMessage = tx.Error.Error()
		response.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := localresponse.NewResponse(true)
	resp.Contacts_Parsed = &rows
	response.JSON(http.StatusOK, resp)
}
