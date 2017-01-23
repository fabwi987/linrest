package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fabwi987/linrest/models"
	"github.com/fabwi987/yaho"
	"github.com/gin-gonic/gin"
)

type Env struct {
	db models.Datastore
}

func main() {
	db, err := models.NewDB("root:trustno1@/test?parseTime=true")
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
	router.GET("/users", env.GetUsersEndpoint)
	router.GET("/users/:id", env.GetSingleUserEndpoint)
	router.GET("/recommendations", env.GetRecommendationsEndpoint)
	router.GET("/recommendations/:id", env.GetRecommendationsByUsersEndpoint)

	router.GET("/test", env.CreateRecommendationsEndpoint)

	router.Run(":" + port)

}

func (env *Env) GetStocksEndpoint(c *gin.Context) {

	stocks, err := env.db.GetStocks()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var symbols string
	for i := 0; i < len(stocks); i++ {
		symbols = symbols + "," + stocks[i].Symbol
	}

	latestStocks, err := yaho.GetStocks(symbols)

	for i := 0; i < len(latestStocks.Query.Results.Quote); i++ {
		var tempfloat float64
		tempfloat, err := strconv.ParseFloat(latestStocks.Query.Results.Quote[i].LastTradePriceOnly, 64)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		stocks[i].LastTradePriceOnly = tempfloat

		tempfloat, err = strconv.ParseFloat(latestStocks.Query.Results.Quote[i].LastTradePriceOnly, 64)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		tempfloat = tempfloat / stocks[i].BuyPrice
		stocks[i].Change = tempfloat
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

func (env *Env) GetUsersEndpoint(c *gin.Context) {

	usrs, err := env.db.GetUsers()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, usrs)
}

func (env *Env) GetSingleUserEndpoint(c *gin.Context) {

	symbol := c.Param("id")
	intid, err := strconv.Atoi(symbol)
	usr, err := env.db.GetSingleUser(intid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, usr)
}

func (env *Env) GetRecommendationsEndpoint(c *gin.Context) {

	recs, err := env.db.GetRecommendations()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, recs)
}

func (env *Env) GetRecommendationsByUsersEndpoint(c *gin.Context) {

	symbol := c.Param("id")
	intid, err := strconv.Atoi(symbol)
	recs, err := env.db.GetRecommendationsByUser(intid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, recs)
}

func (env *Env) CreateRecommendationsEndpoint(c *gin.Context) {

	recs, err := env.db.CreateRecommendation("CBA.ST", 1, 1)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, recs)
}
