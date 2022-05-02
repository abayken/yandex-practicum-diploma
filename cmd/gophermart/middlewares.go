package main

import (
	"net/http"

	"github.com/abayken/yandex-practicum-diploma/internal/creds"
	"github.com/gin-gonic/gin"
)

func SetUserID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")

		if token == "" || err != nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
		} else {
			id, err := creds.Creds{}.Id(token)

			if err != nil {
				ctx.Status(http.StatusInternalServerError)
				ctx.Abort()
			} else {
				ctx.Set("userID", id)
			}
		}

		ctx.Next()
	}
}
