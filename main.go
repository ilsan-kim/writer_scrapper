package main

import (
	"errors"
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
	firstPubDate	string
	novelDesc		string
	genre			string
	writerNickname 	string
	writerEmail		string
}

var baseURL = "https://novel.naver.com/challenge/list?novelId="
var errNoPage = errors.New("no pages to return")
var emailContainer []string

func main() {
	for i:=0 ; i<100 ; i++{
		queryURL := baseURL + strconv.Itoa(i)
		pages := GetPages(queryURL)
		if pages != nil {
			date := GetPubDate(pages)
			novelDesc, _ := GetNovelDesc(pages)
			title := GetNovelTitle(pages)
			genre := GetGenre(pages)
			writerNickname := GetWriterNickname(pages)
			fmt.Printf("%s: %s(%s)[%s] %s. 작가: %s\n", queryURL, title, genre, date, novelDesc, writerNickname)
		}
	}
	fmt.Printf("Found Email : %v\n",len(emailContainer))
}

func GetPages(url string) *goquery.Document {
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
	if checkResp(doc) == nil {
		return doc
	}
	return nil
}

func GetPubDate(doc *goquery.Document) string {
	searchArea := doc.Find(".cont_sub")
	novelList := searchArea.Find("ul")
	date := novelList.Find("li").First().Find(".list_info").Find(".date").Text()
	return date
}

func GetNovelDesc(doc *goquery.Document) (string, string)  {
	var novelDesc = ""
	var writerEmail = ""
	searchArea := doc.Find(".section_area_info")
	searchArea.Each(func(i int, s *goquery.Selection){
		s.Find("p").Each(func(idx int, sel *goquery.Selection) {novelDesc = sel.Text()})
		writerEmail = getEmail(novelDesc)
	})
		return novelDesc, writerEmail
}

func GetNovelTitle(doc *goquery.Document) string {
	title := doc.Find(".book_title").Text()
	return title
}

func GetWriterNickname(doc *goquery.Document) string {
	writerMeta := doc.Find(".writer")
	writerNickname := writerMeta.Find("a").First().Text()
	return writerNickname
}

func GetGenre(doc *goquery.Document) string {
	genre := doc.Find(".genre").Text()
	return genre
}

func getEmail(text string) string {
	re := regexp.MustCompile(`[a-z0-9._%+\-\[]+@[a-z0-9.\-]+\.[a-z\]]{2,4}`)
	is := re.MatchString(text)
	if is == true {
		match := re.FindString(text)
		emailContainer = append(emailContainer, match)
		return match
	}
	return ""
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

func checkResp(doc *goquery.Document) error {
	if len(doc.Text()) == 10001 {
		return errNoPage
	}
	return nil
}