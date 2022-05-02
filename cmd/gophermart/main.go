package main

import (
	"flag"
	"log"

	"github.com/abayken/yandex-practicum-diploma/internal/creds"
	"github.com/abayken/yandex-practicum-diploma/internal/database"
	"github.com/abayken/yandex-practicum-diploma/internal/handlers"
	"github.com/abayken/yandex-practicum-diploma/internal/repositories"
	"github.com/abayken/yandex-practicum-diploma/internal/usecases"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
)

type Config struct {
	RunAddress     string `env:"RUN_ADDRESS" envDefault:":8080"`
	DatabaseURL    string `env:"DATABASE_URI" envDefault:"postgres://abayken:password@localhost:5432/gophermart"`
	AccuralAddress string `env:"ACCURAL_SYSTEM_ADDRESS" envDefault:"localhost:8080"`
}

func main() {
	/// получаем переменные окружения
	cfg := Config{}
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "Адресс сервера")
	flag.StringVar(&cfg.DatabaseURL, "d", cfg.DatabaseURL, "Урл базы данных")
	flag.StringVar(&cfg.AccuralAddress, "r", cfg.AccuralAddress, "Адресс accrual сервиса")

	flag.Parse()

	router := GetRouter(cfg)
	router.Run(cfg.RunAddress)
}

func GetRouter(cfg Config) *gin.Engine {
	router := gin.New()

	storage := database.NewStorage(cfg.DatabaseURL)
	authRepository := repositories.AuthRepository{Storage: storage}
	ordersRepo := repositories.OrdersRepository{Storage: storage}
	withdrawsRepo := repositories.WithdrawRepository{Storage: storage}
	authUseCase := usecases.AuthUseCase{
		Repository:    authRepository,
		Creds:         creds.Creds{},
		OrdersRepo:    ordersRepo,
		WithdrawsRepo: withdrawsRepo,
	}

	ordersUseCase := usecases.OrderUseCase{Repo: ordersRepo}
	handler := handlers.Handler{AuthUseCase: authUseCase, OrdersUseCase: ordersUseCase}

	accrualRepo := repositories.AccrualRepository{BaseURL: cfg.AccuralAddress}
	accrualUseCase := usecases.AccrualUseCase{OrdersRepository: ordersRepo, AccrualRepository: accrualRepo}

	router.POST("/api/user/register", handler.RegisterUser)
	router.POST("/api/user/login", handler.LoginUser)

	authorized := router.Group("/")

	authorized.Use(SetUserID(), ActualizeOrders(accrualUseCase))

	authorized.POST("/api/user/orders", handler.AddOrder)
	authorized.GET("/api/user/orders", handler.Orders)
	authorized.GET("/api/user/balance", handler.Balance)

	return router
}
