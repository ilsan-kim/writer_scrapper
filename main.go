package main

import (
	"github.com/ddhyun93/writerscrapper/scrapper"
	"github.com/labstack/echo"
	"os"
	"strconv"
)

func main() {
	//scrapper.Scrape(1000)
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}

const fileName string = "NovelData.xlsx"

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	num := c.FormValue("num")
	intNum, _ := strconv.Atoi(num)
	scrapper.Scrape(intNum)
	return c.Attachment(fileName, fileName)
}