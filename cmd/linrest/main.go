package main

import (
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/fabwi987/linrest/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Env struct {
	db models.Datastore
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := models.NewDB(os.Getenv("CONNECTION_STRING_LOCAL"))

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
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/stocks", env.GetStocksEndpoint)
	router.GET("/stocks/:symbol", env.GetSingleStocksEndpoint)

	router.GET("/users", env.GetUsersEndpoint)
	router.GET("/users/single/:id", env.GetSingleUserEndpoint)
	router.GET("/users/leaderboard", env.GetUserLeaderboardEndpoint)

	router.GET("/recommendations", env.GetRecommendationsEndpoint)
	router.GET("/recommendations/user/:id", env.GetRecommendationsByUsersEndpoint)
	router.GET("/recommendations/meet/:id", env.GetRecommendationsByMeetEndpoint)

	router.GET("/meets", env.GetMeetsEndpoint)
	router.GET("/meet/single/:id", env.GetSingleMeetEndpoint)
	router.GET("/meet/user/:id", env.GetMeetByUserEndpoint)

	router.GET("/trans/byuser/:id", env.GetTransactionsByUserEndpoint)
	router.GET("/trans/total/:id", env.SumTransactionsByUserEndpoint)

	router.POST("/user", env.CreateUserEndpoint)
	router.POST("/stock", env.CreateStockEndpoint)
	router.POST("/meet", env.CreateMeetEndpoint)
	router.POST("/meet/reward", env.RewardMeetEndpoint)
	//router.POST("/rec", env.CreateRecommendationsEndpoint)
	router.POST("/rec/meet/:id", env.CreateRecommendationsEndpoint)
	router.POST("/trans", env.CreateTransactionEndpoint)

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

	if len(stocks) > 1 {
		latestStocks, err := models.GetStocks(symbols)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
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
			stocks[i].Change = strconv.FormatFloat(tempfloat, 'f', 2, 64)
		}

	} else {
		latestStocks, err := models.GetSingleStocks(symbols)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		var tempfloat float64
		tempfloat, err = strconv.ParseFloat(latestStocks.Query.Results.Quote.LastTradePriceOnly, 64)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		stocks[0].LastTradePriceOnly = tempfloat

		tempfloat, err = strconv.ParseFloat(latestStocks.Query.Results.Quote.LastTradePriceOnly, 64)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		tempfloat = tempfloat / stocks[0].BuyPrice
		stocks[0].Change = strconv.FormatFloat(tempfloat, 'f', 2, 64)
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

func (env *Env) GetUserLeaderboardEndpoint(c *gin.Context) {

	var usr models.Users
	usr, err := env.db.GetUsersLeaderboard()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	sort.Sort(usr)
	c.JSON(200, usr)
}

func (env *Env) GetRecommendationsEndpoint(c *gin.Context) {

	var recs models.Recommendations
	recs, err := env.db.GetRecommendations()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var symbols string
	for i := 0; i < len(recs); i++ {
		symbols = symbols + "," + recs[i].Stck.Symbol
	}

	if len(recs) > 1 {
		latestStocks, err := models.GetStocks(symbols)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		for i := 0; i < len(latestStocks.Query.Results.Quote); i++ {

			tempfloat, err := strconv.ParseFloat(latestStocks.Query.Results.Quote[i].LastTradePriceOnly, 64)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			recs[i].Stck.LastTradePriceOnly = tempfloat

		}

	} else {
		latestStocks, err := models.GetSingleStocks(symbols)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		tempfloat, err := strconv.ParseFloat(latestStocks.Query.Results.Quote.LastTradePriceOnly, 64)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		recs[0].Stck.LastTradePriceOnly = tempfloat
	}

	for i := 0; i < len(recs); i++ {

		dev := ((recs[i].Stck.LastTradePriceOnly / recs[i].Stck.BuyPrice) * 100) - 100
		if dev > 0 {
			recs[i].Stck.Color = "bg-success"
		} else if dev == 0 {
			recs[i].Stck.Color = "bg-warning"
		} else {
			recs[i].Stck.Color = "bg-danger"
		}
		procString := strconv.FormatFloat(dev, 'f', 2, 64)
		subString := procString
		recs[i].Stck.Change = subString
	}

	sort.Sort(recs)

	c.JSON(200, recs)
}

func (env *Env) GetStartEndpoint(c *gin.Context) {
	c.HTML(http.StatusOK, "start.html", nil)
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

func (env *Env) GetRecommendationsByMeetEndpoint(c *gin.Context) {

	symbol := c.Param("id")
	intid, err := strconv.Atoi(symbol)
	recs, err := env.db.GetRecommendationsByMeet(intid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, recs)
}

func (env *Env) CreateRecommendationsEndpoint(c *gin.Context) {
	symbol := c.Param("id")
	meetid, err := strconv.Atoi(symbol)

	user, err := strconv.Atoi(c.PostForm("user"))
	stcken, err := models.GetSingleStocks(c.PostForm("symbol"))
	lasttradeprice, err := strconv.ParseFloat(stcken.Query.Results.Quote.LastTradePriceOnly, 64)

	_, err = env.db.CreateStock(stcken.Query.Results.Quote.Symbol, stcken.Query.Created, lasttradeprice, 0, 0, stcken.Query.Results.Quote.Name, lasttradeprice)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	recs, err := env.db.CreateRecommendation(stcken.Query.Results.Quote.Symbol, user, meetid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	/**
	recuser, err := env.db.GetSingleUser(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	//err = models.SendRecommendationMail(recuser.Mail, recuser.Name)
	//err = models.SendRecommendationText(recuser.Phone, recuser.Name, stcken.Query.Results.Quote.Name)*/

	c.JSON(201, recs)
}

func (env *Env) CreateUserEndpoint(c *gin.Context) {

	usr, err := env.db.CreateUser(c.Query("name"), c.Query("phone"), c.Query("mail"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, usr)
}

func (env *Env) CreateStockEndpoint(c *gin.Context) {

	stcken, err := models.GetSingleStocks(c.Query("symbol"))

	buyprice, err := strconv.ParseFloat(c.Query("buyprice"), 64)
	numberofshares, err := strconv.Atoi(c.Query("numberofshares"))
	lasttradeprice, err := strconv.ParseFloat(stcken.Query.Results.Quote.LastTradePriceOnly, 64)

	log.Println("Hämtar en stock:")
	log.Println(stcken.Query.Results.Quote.Name)
	log.Println(stcken.Query.Results.Quote.LastTradePriceOnly)

	stck, err := env.db.CreateStock(stcken.Query.Results.Quote.Symbol, stcken.Query.Created, buyprice, numberofshares, 0, stcken.Query.Results.Quote.Name, lasttradeprice)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, stck)
}

func (env *Env) CreateMeetEndpoint(c *gin.Context) {
	user, err := strconv.Atoi(c.PostForm("user"))
	met, err := env.db.CreateMeet(c.PostForm("location"), c.PostForm("date"), c.PostForm("text"), user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(201, met)
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

func (env *Env) GetMeetByUserEndpoint(c *gin.Context) {

	symbol := c.Param("id")
	intid, err := strconv.Atoi(symbol)
	meet, err := env.db.GetMeetsByUser(intid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, meet)
}

func (env *Env) GetTransactionsByUserEndpoint(c *gin.Context) {

	symbol := c.Param("id")
	intid, err := strconv.Atoi(symbol)
	trans, err := env.db.GetTransactionsByUser(intid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, trans)
}

func (env *Env) SumTransactionsByUserEndpoint(c *gin.Context) {

	symbol := c.Param("id")
	intid, err := strconv.Atoi(symbol)
	trans := env.db.SumTransactionsByUser(intid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, trans)
}

func (env *Env) CreateTransactionEndpoint(c *gin.Context) {

	rec, err := strconv.Atoi(c.Query("recommendation"))
	user, err := strconv.Atoi(c.Query("user"))
	reward, err := strconv.Atoi(c.Query("reward"))

	usr, err := env.db.CreateTransaction(rec, user, reward)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(200, usr)
}

func (env *Env) RewardMeetEndpoint(c *gin.Context) {

	symbol := c.Query("id")
	intid, err := strconv.Atoi(symbol)
	var recs models.Recommendations
	recs, err = env.db.GetRecommendationsByMeet(intid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var symbols string
	for i := 0; i < len(recs); i++ {
		symbols = symbols + "," + recs[i].Stck.Symbol
	}

	latestStocks, err := models.GetStocks(symbols)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	for i := 0; i < len(latestStocks.Query.Results.Quote); i++ {

		tempfloat, err := strconv.ParseFloat(latestStocks.Query.Results.Quote[i].LastTradePriceOnly, 64)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		recs[i].Stck.LastTradePriceOnly = tempfloat

	}

	for i := 0; i < len(recs); i++ {
		dev := ((recs[i].Stck.LastTradePriceOnly / recs[i].Stck.BuyPrice) * 100) - 100
		procString := strconv.FormatFloat(dev, 'f', 2, 64)
		subString := procString
		recs[i].Stck.Change = subString
	}

	sort.Sort(recs)

	resp, err := env.db.CreateTransaction(recs[1].ID, recs[1].Usr.ID, 100)

	c.JSON(200, resp)
}
