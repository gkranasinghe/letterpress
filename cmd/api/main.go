package main

import (
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"

	"github.com/gin-gonic/gin"
	"github.com/gkranasinghe/letterpress/db"
	"github.com/gkranasinghe/letterpress/handler"
	"github.com/rs/zerolog"
)

func main() {
	var dbPort int
	var err error
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	port := os.Getenv("POSTGRES_PORT")
	if dbPort, err = strconv.Atoi(port); err != nil {
		logger.Err(err).Msg("failed to parse database port")
		os.Exit(1)
	}
	dbConfig := db.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     dbPort,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
		Logger: logger,
	}
	logger.Info().Interface("config", &dbConfig).Msg("config:")
	dbInstance, err := db.Init(dbConfig)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}
	logger.Info().Msg("Database connection established")
    
	cfg := elasticsearch.Config{
		Addresses: []string{
		  "http://elasticsearch:9200",
		},
		// ...
	  }


	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}
	
	res, err := esClient.Info()
	if err != nil {
		logger.Err(err).Msg("Error getting response")
		os.Exit(1)
	 
	}
	
	defer res.Body.Close()

	// esClient, err := elasticsearch.NewDefaultClient()
	// if err != nil {
	// 	logger.Err(err).Msg("Connection failed")
	// 	os.Exit(1)
	// }

	h := handler.New(dbInstance, esClient, logger)
	router := gin.Default()
	rg := router.Group("/v1")
	h.Register(rg)
	router.Run(":8080")
}
