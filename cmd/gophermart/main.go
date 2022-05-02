package main

import (
	"flag"
	"log"

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
	AccuralAddress string `env:"ACCURAL_SYSTEM_ADDRESS"`
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

	flag.Parse()

	router := GetRouter(cfg)
	router.Run(cfg.RunAddress)
}

func GetRouter(cfg Config) *gin.Engine {
	router := gin.New()

	storage := database.NewStorage(cfg.DatabaseURL)
	authRepository := repositories.AuthRepository{Storage: storage}
	authUseCase := usecases.AuthUseCase{Repository: authRepository}

	ordersRepo := repositories.OrdersRepository{Storage: storage}
	ordersUseCase := usecases.OrderUseCase{Repo: ordersRepo}
	handler := handlers.Handler{AuthUseCase: authUseCase, OrdersUseCase: ordersUseCase}

	router.POST("/api/user/register", handler.RegisterUser)
	router.POST("/api/user/login", handler.LoginUser)

	authorized := router.Group("/")

	authorized.Use(SetUserID())

	authorized.POST("/api/user/orders", handler.AddOrder)

	return router
}
