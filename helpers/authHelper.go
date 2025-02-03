package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func AuthorizeRole(c *gin.Context, userId string) error {
	userType := c.GetString("userType")
	uid := c.GetString("uid")

	// grant access to Admin
	if userType == "ADMIN" {
		return nil
	}

	// grant access to user if the user id matches
	if userType == "USER" && uid == userId {
		return nil
	}

	// deny all other permissions
	return errors.New("unauthorized request")
}
