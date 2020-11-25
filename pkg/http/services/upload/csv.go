package upload

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/dathan/go-web-backend/pkg/entities"
	localresponse "github.com/dathan/go-web-backend/pkg/http/response"
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

	err = csvUp.saveUploadToDatabase(reader, user)
	if err != nil {
		resp := localresponse.NewResponse(false)
		resp.ErrorMessage = fmt.Sprintf("validationError: %s", err.Error())
		response.JSON(http.StatusNotAcceptable, resp)
		return

	}

	resp := localresponse.NewResponse(true)
	response.JSON(http.StatusOK, resp)
	return

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
func (c *csvUp) saveUploadToDatabase(reader *multipart.Reader, user *entities.User) error {
	var files []string

	//copy each part to destination.
	for {
		row := &entities.Contacts_Paths{}
		row.OwnerID = user.ID
		part, err := reader.NextPart()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
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
			return err
		}

		defer func() {
			if err := dst.Close(); err != nil {
				fmt.Errorf("Close Error: %s\n", err)
			}

		}()
		byte_data, err := ioutil.ReadAll(part)
		if err != nil {
			return err
		}
		row.CSVData = string(byte_data)
		fmt.Printf("What Is row.Data\n:%s\n%v\n", string(byte_data), byte_data)

		if _, err := io.Copy(dst, part); err != nil {
			return err
		}

		tx := database.GetConnection().Create(row)
		if tx.Error != nil {
			return tx.Error
		}

	}

	fmt.Printf("Uploaded %d files\n\t%+v\n", len(files), files)

	return nil
}

func (c *csvUp) getAllParsedFiles(user *entities.User) {

}
