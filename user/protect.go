package user

import (
	"github.com/gin-gonic/gin"

	e "github.com/eirka/eirka-libs/errors"
)

// Protect will check to see if a user has the correct permissions
func Protect() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get parameters from validate middleware
		params := c.MustGet("params").([]uint)

		// get userdata from session middleware
		userdata := c.MustGet("userdata").(User)

		// check if user is authorized
		if !userdata.IsAuthorized(params[0]) {
			c.JSON(e.ErrorMessage(e.ErrForbidden))
			c.Error(e.ErrForbidden).SetMeta("user.Protect.IsAuthorized")
			c.Abort()
			return
		}

		// this route was protected
		c.Set("protected", true)

		c.Next()

	}
}
