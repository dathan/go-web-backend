package entities

import userentity "github.com/dathan/go-web-backend/pkg/entities/user"

//CommonResponse puts a envelope around the response of everything so it's common
type CommonResponse struct {
	OK           bool             `json:"ok"`
	ErrorMessage string           `json:"errorMsg,omitempty"`
	User         *userentity.User `json:"user,omitempty"`
	Token        string           `json:"token,omitempty"`
	RefreshToken string           `json:"refresh_token,omitempty"`
	// add more entities below
}

func NewResponse(status bool) *CommonResponse {
	return &CommonResponse{
		OK: status,
	}
}
