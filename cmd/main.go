package main

import (
	"log"

	"github.com/TOPfiIT/vacancy-service/internal/config"
	"github.com/TOPfiIT/vacancy-service/internal/db"
	"github.com/TOPfiIT/vacancy-service/internal/http/handlers"
	"github.com/TOPfiIT/vacancy-service/internal/services"
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

	//create routes
	r.POST("/vacancies", vacancyHandler.CreateVacancy)
	r.GET("/vacancies/company/:company_id", vacancyHandler.GetVacancies)
	r.GET("/vacancies/:vacancy_id", vacancyHandler.GetVacancy)
	r.GET("/vacancies/:vacancy_id/interview", vacancyHandler.GetInterviewResults)
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
