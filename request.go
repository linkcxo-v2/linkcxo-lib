package linkcxo

import (
	"errors"

	"github.com/labstack/echo/v4"
)

const (
	//ContentCredKey -
	ContentCredKey = " credentials"
)

// RequestUtils -
type RequestUtils struct {
}

// GetCredential -
func (ru RequestUtils) GetCredential(c echo.Context) *UserCredential {
	var cred = c.Get(ContentCredKey)
	if cred == nil {
		return nil
	}
	uc := cred.(UserCredential)
	return &uc
}

// SetCredential -
func (ru RequestUtils) SetCredential(c echo.Context, userCred UserCredential) {
	c.Set(ContentCredKey, userCred)
}

// GetUserID -
func (ru RequestUtils) GetUserID(c echo.Context) (string, error) {
	userCred := ru.GetCredential(c)
	if userCred != nil && userCred.UserID != "" {
		return userCred.UserID, nil
	}
	return "", errors.New("access forbidden")
}
