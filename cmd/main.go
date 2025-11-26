package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
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
		PublicKey: getPublicKey(),
	}
	authMiddleware := middlewares.AuthMiddleware(authConfig)

	//create routes
	protected := r.Group("/")
	protected.Use(authMiddleware)
	{
		protected.POST("/vacancies", vacancyHandler.CreateVacancy)
		protected.GET("/vacancies/company", vacancyHandler.GetVacancies)
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

func getPublicKey() *ecdsa.PublicKey {
	pemKey := `-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEQy0TSI95s6hkPCFyicGS8oVoAzB7\nyiQhjrsy9KV4QivobyYsb89YxM7BJ+HHGhL83+DcKwRViLm3SlMV6NP/IQ==\n-----END PUBLIC KEY-----`

	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		log.Fatal("Failed to parse PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Not an ECDSA public key")
	}

	return ecdsaPubKey
}
