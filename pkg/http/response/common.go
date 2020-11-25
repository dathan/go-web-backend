package response

import "github.com/dathan/go-web-backend/pkg/entities"

//CommonResponse puts a envelope around the response of everything so it's common
type CommonResponse struct {
	OK              bool                        `json:"ok"`
	ErrorMessage    string                      `json:"errorMsg,omitempty"`
	User            *entities.User              `json:"user,omitempty"`
	Token           string                      `json:"token,omitempty"`
	RefreshToken    string                      `json:"refresh_token,omitempty"`
	Contacts_Paths  *[]entities.Contacts_Paths  `json:"contacts_paths,omitempty"`
	Contacts_Parsed *[]entities.Contacts_Parsed `json:"contacts_parsed,omitempty"`
	// add more entities below
}

func NewResponse(status bool) *CommonResponse {
	return &CommonResponse{
		OK: status,
	}
}
