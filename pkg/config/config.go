package config

import (
	"fmt"
	"log"
	"os"

	"github.com/AgungAryansyah/filkompedia-be-insecure/internal/handler/rest"
	"github.com/AgungAryansyah/filkompedia-be-insecure/internal/repository"
	"github.com/AgungAryansyah/filkompedia-be-insecure/internal/service"
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/bcrypt"
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/jwt"
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/logger"
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/middleware"
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/midtrans"
	monitoring "github.com/AgungAryansyah/filkompedia-be-insecure/pkg/prometheus"
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/smtp"
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/supabase"
	val "github.com/AgungAryansyah/filkompedia-be-insecure/pkg/validator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type Config struct {
	DB    *sqlx.DB
	Redis *redis.Client
	App   *fiber.App
}

func LoadEnv() {
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func StartUp(config *Config) {
	bcrypt := bcrypt.Init()
	jwt := jwt.Init()
	smtp := smtp.LoadSMTPCredentials()
	midtrans := midtrans.NewMidtrans()
	promMetrics := monitoring.Start()
	logrus := logger.SetupLogger()
	supabase := supabase.New()

	validator := validator.New()
	val.RegisterValidator(validator)

	repository := repository.NewRepository(config.DB, config.Redis)
	service := service.NewService(repository, bcrypt, jwt, smtp, midtrans, supabase)

	middleware := middleware.Init(jwt, service, promMetrics, logrus)

	config.App.Use(middleware.PromMiddleware)
	config.App.Use(middleware.LogrusMiddleware)

	rest := rest.NewRest(config.App, service, middleware, validator)
	rest.RegisterRoutes()

	rest.Start(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
