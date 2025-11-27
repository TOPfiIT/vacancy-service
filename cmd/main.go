package main

import (
	"log"

	"github.com/TOPfiIT/vacancy-service/internal/config"
	"github.com/TOPfiIT/vacancy-service/internal/db"
	"github.com/TOPfiIT/vacancy-service/internal/http/handlers"
	"github.com/TOPfiIT/vacancy-service/internal/services"
	"github.com/TOPfiIT/vacancy-service/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	//init cfg
	cfg := config.MustLoad()

	publicKey, err := cfg.GetPublicKey()
	if err != nil {
		log.Fatalf("Failed to load public key: %v", err)
	}
	log.Println("✓ Public key loaded")

	// init postgres
	repository := db.InitPostgres(cfg)
	// defer repository.Close()

	//init services
	vacancyService := services.NewVacancyService(repository)

	//init handlers
	vacancyHandler := handlers.NewVacancyHandler(vacancyService)

	//creater router
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	authConfig := middlewares.AuthConfig{
		PublicKey: publicKey,
	}
	authMiddleware := middlewares.AuthMiddleware(authConfig)

	//create routes
	protected := r.Group("/")
	protected.Use(authMiddleware)
	{
		protected.POST("/vacancies", vacancyHandler.CreateVacancy)
		protected.GET("/vacancies/company", vacancyHandler.GetVacancies)
		protected.GET("/vacancies/company/:vacancy_id", vacancyHandler.GetVacancyFront)
		protected.GET("/vacancies/:vacancy_id/interview", vacancyHandler.GetInterviewResults)
	}

	r.GET("/vacancies/:vacancy_id", vacancyHandler.GetVacancy)
	r.POST("/vacancies/:vacancy_id/interview", vacancyHandler.AddInterviewResult)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	//start server
	log.Printf("Service running on :%d", cfg.Server.Port)
	if err := r.Run(":8086"); err != nil {
		log.Fatal(err)
	}
}
