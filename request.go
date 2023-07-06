package linkcxo

import (
	"strconv"

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
func (ru RequestUtils) GetCredential(c echo.Context) UserCredential {
	return c.Get(ContentCredKey).(UserCredential)
}

// SetCredential -
func (ru RequestUtils) SetCredential(c echo.Context, userCred UserCredential) {
	c.Set(ContentCredKey, userCred)
}

// GetUserID -
func (ru RequestUtils) GetUserID(c echo.Context) (int64, error) {
	userCred := ru.GetCredential(c)
	return strconv.ParseInt(userCred.UserID, 10, 0)
}
