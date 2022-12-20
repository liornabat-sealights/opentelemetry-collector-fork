package slauth

import "fmt"

type ErrorAuthenticationResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *ErrorAuthenticationResponse) Validate() error {
	if r.Code == 0 {
		return fmt.Errorf("code has 0 value")
	}
	return nil
}

func NewErrorAuthenticationResponse() *ErrorAuthenticationResponse {
	return &ErrorAuthenticationResponse{
		Code:    0,
		Message: "",
	}
}

func (r *ErrorAuthenticationResponse) SetCode(code int) *ErrorAuthenticationResponse {
	r.Code = code
	return r
}

func (r *ErrorAuthenticationResponse) SetMessage(message string) *ErrorAuthenticationResponse {
	r.Message = message
	return r
}

var _ model = (*ErrorAuthenticationResponse)(nil)
