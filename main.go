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
var emailContainer []string

func main() {
	for i:=0 ; i<10000 ; i++{
		queryURL := baseURL + strconv.Itoa(i)
		resp := GetPages(queryURL)
		fmt.Println(queryURL, ": ", resp)
		email := getEmail(resp)
		if email != "" {
			emailContainer = append(emailContainer, email)
		}
	}
	fmt.Printf("Found Email : %v\n",len(emailContainer))
}

func GetPages(url string) string {
	var novelDesc = ""
	var writerEmail = ""
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
	is := re.MatchString(text)
	if is == true {
		match := re.FindString(text)
		return match
	}
	return ""
}
