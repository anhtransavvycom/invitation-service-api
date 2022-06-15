package middleware

import (
	"app-invite-service/common"
	"app-invite-service/component"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// Recover needs AppContext due to 2 reasons below
// - Log error to DB
// - Get environment settings (prod/stag/dev)
func Recover(_ component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				// if error is an AppError
				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					// Gin has its own `Recover`, that wraps our `Recover`
					// Gin can dumb your error to the terminal when we call `panic` here.
					// It makes dev easier to trace bugs.
					// Call `return` here won't dumb error in the terminal
					panic(err)
				}

				appErr := common.ErrInternal(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err)
			}
		}()

		c.Next()
	}
}

func FiberRecover(_ component.AppContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

				// if error is an AppError
				if appErr, ok := err.(*common.AppError); ok {
					c.Status(appErr.StatusCode).JSON(appErr)
					panic(err)
				}

				appErr := common.ErrInternal(err.(error))
				c.Status(appErr.StatusCode).JSON(appErr)
				//c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err)
			}
		}()

		c.Next()

		return nil
	}
}
