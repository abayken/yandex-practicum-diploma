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
	AccrualUseCase  usecases.AccrualUseCase
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

	orderNumberAsString := string(body)

	orderNumber, err := strconv.Atoi(orderNumberAsString)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	userID := ctx.GetInt("userID")

	added, err := handler.OrdersUseCase.Add(userID, orderNumber)

	if added {
		// go func() {
		// 	handler.AccrualUseCase.ActualizeOrders(orderNumberAsString)
		// }()

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

	if len(orders) == 0 {
		ctx.Status(http.StatusNoContent)

		return
	}

	type OrderView struct {
		Number     string  `json:"number"`
		Status     string  `json:"status"`
		Accural    float32 `json:"accural,omitempty"`
		UploadedAt string  `json:"uploaded_at"`
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
			Accural:    float32(order.Accrual) / 100,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (handler *Handler) Balance(ctx *gin.Context) {
	type BalanceView struct {
		Current   float32 `json:"current"`
		Withdrawn float32 `json:"withdrawn"`
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
		Sum   float32 `json:"sum"`
	}

	var request Request

	if err := ctx.BindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	userID := ctx.GetInt("userID")
	err := handler.WithdrawUseCase.Withdraw(userID, request.Order, request.Sum)

	if err != nil {
		var invalidOrderNumberError *custom_errors.InvalidOrderNumber
		var insufficientFundsError *custom_errors.InsufficientFundsError

		if errors.As(err, &invalidOrderNumberError) {
			ctx.Status(http.StatusUnprocessableEntity)
		} else if errors.As(err, &insufficientFundsError) {
			ctx.Status(http.StatusPaymentRequired)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	ctx.Status(http.StatusOK)
}

func (hander *Handler) FakeAccural(ctx *gin.Context) {
	type Response struct {
		Order   string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float32 `json:"accrual,omitempty"`
	}

	ctx.JSON(http.StatusOK, Response{Order: ctx.Param("number"), Status: "INVALID"})
}

func (handler *Handler) Withdrawals(ctx *gin.Context) {
	userID := ctx.GetInt("userID")

	withdrawals, err := handler.WithdrawUseCase.Withdrawals(userID)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	if len(withdrawals) == 0 {
		ctx.Status(http.StatusNoContent)

		return
	}

	type WithdrawView struct {
		Order       string  `json:"order"`
		Sum         float32 `json:"sum"`
		ProcessedAt string  `json:"processed_at"`
	}

	var withdrawalsList []WithdrawView

	for _, withdraw := range withdrawals {
		withdrawView := WithdrawView{
			Order:       withdraw.OrderNumber,
			Sum:         float32(withdraw.Sum) / 100,
			ProcessedAt: withdraw.AddedAt.Time.Format(time.RFC3339),
		}

		withdrawalsList = append(withdrawalsList, withdrawView)
	}

	ctx.JSON(http.StatusOK, withdrawalsList)
}
