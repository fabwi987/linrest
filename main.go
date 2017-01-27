package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"time"

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
	router.GET("/meet", env.GetMeetsEndpoint)
	router.GET("/meet/:id", env.GetSingleMeetEndpoint)

	router.POST("/user", env.CreateUserEndpoint)
	router.POST("/stock", env.CreateStockEndpoint)
	router.POST("/meet", env.CreateMeetEndpoint)
	router.POST("/rec", env.CreateRecommendationsEndpoint)

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

	user, err := strconv.Atoi(c.Query("iduser"))
	meet, err := strconv.Atoi(c.Query("idmeet"))

	recs, err := env.db.CreateRecommendation(c.Query("symbol"), user, meet)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, recs)
}

func (env *Env) CreateUserEndpoint(c *gin.Context) {

	usr, err := env.db.CreateUser(c.Query("name"), c.Query("phone"), c.Query("mail"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, usr)
}

func (env *Env) CreateStockEndpoint(c *gin.Context) {

	stcken, err := yaho.GetSingleStocks(c.Query("symbol"))

	buyprice, err := strconv.ParseFloat(c.Query("buyprice"), 64)
	numberofshares, err := strconv.Atoi(c.Query("numberofshares"))
	lasttradeprice, err := strconv.ParseFloat(stcken.Query.Results.Quote.LastTradePriceOnly, 64)

	stck, err := env.db.CreateStock(stcken.Query.Results.Quote.Symbol, stcken.Query.Created, buyprice, numberofshares, 0, stcken.Query.Results.Quote.Name, lasttradeprice)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, stck)
}

func (env *Env) CreateMeetEndpoint(c *gin.Context) {

	met, err := env.db.CreateMeet(c.Query("location"), time.Now(), c.Query("text"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, met)
}

func (env *Env) GetMeetsEndpoint(c *gin.Context) {

	meets, err := env.db.GetMeets()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, meets)
}

func (env *Env) GetSingleMeetEndpoint(c *gin.Context) {

	symbol := c.Param("id")
	intid, err := strconv.Atoi(symbol)
	meet, err := env.db.GetSingleMeet(intid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, meet)
}
