package response

import "encoding/json"

const BadRoleProvidedMsg = "bad role provided"

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func ToErrorResponse(e string, code int) ([]byte, error) {
	r := ErrorResponse{
		e,
		code,
	}
	return json.Marshal(&r)
}
