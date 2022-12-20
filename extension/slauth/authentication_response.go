package slauth

import (
	"encoding/json"
	"fmt"
)

type AuthenticationResponse struct {
	CustomerId   string `json:"customerId"`
	ValidUntil   int64  `json:"validUntil"`
	ResponseInfo *HttpResponse
	ErrorInfo    *ErrorAuthenticationResponse
}

func (r *AuthenticationResponse) Validate() error {
	if r.CustomerId == "" {
		return fmt.Errorf("AUTH RESPONSE STRUCT:customerId is empty")
	}
	return nil
}

func NewAuthenticationResponse() *AuthenticationResponse {
	return &AuthenticationResponse{
		CustomerId:   "",
		ResponseInfo: NewHttpResponse(),
		ErrorInfo:    NewErrorAuthenticationResponse(),
	}
}

func (r *AuthenticationResponse) SetValidUntil(validUntil int64) *AuthenticationResponse {
	r.ValidUntil = validUntil
	return r
}

func (r *AuthenticationResponse) SetCustomerId(customerId string) *AuthenticationResponse {
	r.CustomerId = customerId
	return r
}

func (r *AuthenticationResponse) constructError() (authErr error) {
	if authErr = json.Unmarshal(r.ResponseInfo.Body, &(r.ErrorInfo)); authErr != nil {
		return authErr
	}

	if authErr = r.ErrorInfo.Validate(); authErr != nil {
		return authErr
	}

	return fmt.Errorf(fmt.Sprintf("ResponseCode: %d, Error Message: %s, Code: %d ", r.ResponseInfo.StatusCode, r.ErrorInfo.Message, r.ErrorInfo.Code))
}

var _ model = (*AuthenticationResponse)(nil)
