package slauth

import "fmt"

type AuthCacheData struct {
	token      string
	customerId string
}

func NewAuthCacheData() *AuthCacheData {
	return &AuthCacheData{
		token:      "",
		customerId: "",
	}
}

func (d *AuthCacheData) Validate() error {
	if d.token == "" {
		return fmt.Errorf("CACHE: token is empty")
	}
	if d.customerId == "" {
		return fmt.Errorf("CACHE: customerId is empty")
	}
	return nil
}

func (d *AuthCacheData) setToken(token string) *AuthCacheData {
	d.token = token
	return d
}

func (d *AuthCacheData) setCustomerId(customerId string) *AuthCacheData {
	d.customerId = customerId
	return d
}

var _ model = (*AuthCacheData)(nil)
