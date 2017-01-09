package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fabwi987/linrest/models"
	"github.com/gin-gonic/gin"
)

type Env struct {
	db models.Datastore
}

func main() {
	db, err := models.NewDB("root:trustno1@/test")
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/stocks", env.GetStocksEndpoint)
	router.GET("/stocks/:symbol", env.GetSingleStocksEndpoint)

	router.Run(":" + port)

}

func (env *Env) GetStocksEndpoint(c *gin.Context) {

	stocks, err := env.db.GetStocks()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, stocks)
}

func (env *Env) GetSingleStocksEndpoint(c *gin.Context) {

	symbol := c.Param("symbol")
	stock, err := env.db.GetSingleStock(symbol)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, stock)
}
