package validate

import (
	"github.com/gin-gonic/gin"

	e "github.com/eirka/eirka-libs/errors"
)

// ValidateParams will loop through the route parameters to make sure theyre uint
func ValidateParams() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Params != nil {

			var params []uint

			for _, param := range c.Params {

				pid, err := ValidateParam(param.Value)
				if err != nil {
					c.JSON(e.ErrorMessage(e.ErrInvalidParam))
					c.Error(err).SetMeta("validate.ValidateParams")
					c.Abort()
					return
				}

				params = append(params, pid)

			}

			c.Set("params", params)

		}

		c.Next()

	}
}
