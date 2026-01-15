package middlewares

import (
	"log"
	"net/http"
	"pdf-generator/internal/domain"

	"github.com/gin-gonic/gin"
)

func Errorhandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next() // Process request first

		// Check if there are any errors
		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			// Check if it's our custom error (handles HttpError and all subtypes)
			if httpErr, ok := err.(domain.AppErrorInterface); ok {
				ctx.JSON(httpErr.GetCode(), gin.H{"message": httpErr.GetMessage()})
				return
			}

			log.Fatalf("Error: %v", err)

			// Default to 500 for unknown errors
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong, please try again",
			})
		}
	}
}
