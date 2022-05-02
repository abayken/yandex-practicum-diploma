package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/abayken/yandex-practicum-diploma/internal/custom_errors"
	"github.com/abayken/yandex-practicum-diploma/internal/usecases"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	AuthUseCase     usecases.AuthUseCase
	OrdersUseCase   usecases.OrderUseCase
	WithdrawUseCase usecases.WithdrawUseCase
}

func (handler *Handler) RegisterUser(ctx *gin.Context) {
	type Request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var request Request

	if err := ctx.BindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	token, err := handler.AuthUseCase.Register(request.Login, request.Password)

	if err != nil {
		var userExistsError *custom_errors.AlreadyExistsUserError
		if errors.As(err, &userExistsError) {
			ctx.Status(http.StatusConflict)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
	ctx.Status(http.StatusOK)
}

func (handler *Handler) LoginUser(ctx *gin.Context) {
	type Request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var request Request

	if err := ctx.BindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	token, err := handler.AuthUseCase.Login(request.Login, request.Password)

	if err != nil {
		var invalidCredsError *custom_errors.InvalidCredentialsError
		if errors.As(err, &invalidCredsError) {
			ctx.Status(http.StatusUnauthorized)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
	ctx.Status(http.StatusOK)
}

/// Обработчик /api/user/orders
func (handler *Handler) AddOrder(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	orderNumber, err := strconv.Atoi(string(body))

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	userID := ctx.GetInt("userID")

	added, err := handler.OrdersUseCase.Add(userID, orderNumber)

	if added {
		ctx.Status(http.StatusAccepted)

		return
	}

	if err != nil {
		var invalidOrderNumberError *custom_errors.InvalidOrderNumber
		var orderAlreadyAddedError *custom_errors.OrderAlreadyAddedError

		if errors.As(err, &invalidOrderNumberError) {
			ctx.Status(http.StatusUnprocessableEntity)

			return
		}

		if errors.As(err, &orderAlreadyAddedError) {
			if orderAlreadyAddedError.UserID == userID {
				ctx.Status(http.StatusOK)
			} else {
				ctx.Status(http.StatusConflict)
			}

			return
		}

		ctx.Status(http.StatusInternalServerError)
	}
}

func (handler *Handler) Orders(ctx *gin.Context) {
	userID := ctx.GetInt("userID")
	orders, err := handler.OrdersUseCase.GetOrders(userID)

	type OrderView struct {
		Number     string `json:"number"`
		Status     string `json:"status"`
		Accural    string `json:"accural,omitempty"`
		UploadedAt string `json:"uploaded_at"`
	}

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	var response []OrderView

	for _, order := range orders {
		response = append(response, OrderView{
			Number:     order.Number,
			Status:     order.Status,
			UploadedAt: order.AddedAt.Time.Format(time.RFC3339),
		})
	}

	if response == nil {
		ctx.Status(http.StatusNoContent)
	} else {
		ctx.JSON(http.StatusOK, response)
	}
}

func (handler *Handler) Balance(ctx *gin.Context) {
	type BalanceView struct {
		Current   int `json:"current"`
		Withdrawn int `json:"withdrawn"`
	}

	userID := ctx.GetInt("userID")

	balance, err := handler.AuthUseCase.GetBalance(userID)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, BalanceView{Current: balance.Current, Withdrawn: balance.TotalWithdrawn})
}

func (handler *Handler) Withdraw(ctx *gin.Context) {
	type Request struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	var request Request

	if err := ctx.BindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	//userID := ctx.GetInt("userID")

}
