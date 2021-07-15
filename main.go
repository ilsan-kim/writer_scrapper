package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type novelData struct {
	url				string
	title			string
	lastUpdate		string
	novelDesc		string
	writerNickname 	string
	writerEmail		string
}

var baseURL = "https://novel.naver.com/challenge/list?novelId="

func main() {
	for i:=0 ; i<50 ; i++{
		queryURL := baseURL + strconv.Itoa(i)
		fmt.Println(queryURL, ": ",GetPages(queryURL))
	}
}

func GetPages(url string) (string) {
	var novelDesc string = ""
	var writerEmail string = ""
	res, err := http.Get(url)
	checkErr(url, err)
	checkStatus(url, res.StatusCode)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(url, err)

	searchArea := doc.Find(".section_area_info")
	searchArea.Each(func(i int, s *goquery.Selection){
		s.Find("p").Each(func(idx int, sel *goquery.Selection) {
			novelDesc = sel.Text()
		})
		writerEmail = getEmail(novelDesc)
		fmt.Println(writerEmail)
	})
	//doc.Find(".section_area_info").Each(func(i int, s *goquery.Selection) {
	//	s.Find("p").Each(func(idx int, sel *goquery.Selection) {
	//		novelDesc = sel.Text()
	//	})
	//})
	return novelDesc
}


func checkStatus(url string, code int) {
	if code != 200 {
		fmt.Printf("URL: %s Return Status Code %v\n", url, code)
	}
}

func checkErr(url string, err error) {
	if err != nil {
		log.Fatalf("URL: %s Return Err %s\n", url, err)
	}
}

func getEmail(text string) string {
	re := regexp.MustCompile(`[a-z0-9._%+\-\[]+@[a-z0-9.\-]+\.[a-z\]]{2,4}`)
	match := re.FindString(text)
	return match
}
