package main

import (
	"fmt"
	"net/http"

	"github.com/abayken/yandex-practicum-diploma/internal/creds"
	"github.com/abayken/yandex-practicum-diploma/internal/usecases"
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

func ActualizeOrders(usecase usecases.AccrualUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetInt("userID")
		err := usecase.ActualizeOrders(userID)

		if err != nil {
			fmt.Println("error occured while updating data")
			ctx.Status(http.StatusInternalServerError)

			return
		}

		ctx.Next()
	}
}
