//package practice
//
//import (
//	"errors"
//	"fmt"
//	"github.com/PuerkitoBio/goquery"
//	"net/http"
//	"strconv"
//	"time"
//)
//
//var baseURL = "https://novel.naver.com/challenge/list?novelId="
//var errNoPage = errors.New("no pages to return")
//
//func main() {
//	c := make(chan *goquery.Document)
//	for i:=1 ; i<30000 ; i++ {
//		queryURL := baseURL + strconv.Itoa(i)
//		//GetPagesNormal(queryURL)
//		go GetPages(queryURL, c)
//		time.Sleep(time.Microsecond *5500)
//	}
//	for i:=1 ; i<30000 ; i++ {
//		novelData := <- c
//		fmt.Println(novelData)
//	}
//}
//
//
//func GetPages(url string, c chan <- *goquery.Document)  {
//	fmt.Println(url)
//	res, err := http.Get(url)
//	checkErr(url, err)
//	checkStatus(url, res.StatusCode)
//	defer res.Body.Close()
//
//	doc, err := goquery.NewDocumentFromReader(res.Body)
//	checkErr(url, err)
//
//	if checkResp(doc) != nil {
//		c <- nil
//	}
//	c <- doc
//}
//
//func GetPagesNormal(url string) *goquery.Document {
//	fmt.Println(url)
//	res, err := http.Get(url)
//	checkErr(url, err)
//	checkStatus(url, res.StatusCode)
//	defer res.Body.Close()
//
//	doc, err := goquery.NewDocumentFromReader(res.Body)
//	checkErr(url, err)
//
//	if checkResp(doc) != nil {
//		return nil
//	}
//	return doc
//}
//
//func checkErr(url string, err error) {
//	if err != nil {
//		fmt.Printf("URL: %s Return Err %s\n", url, err)
//	}
//}
//
//func checkStatus(url string, code int) {
//	if code != 200 {
//		fmt.Printf("URL %s Return Status Code %v\n", url, code)
//	}
//}
//
//func checkResp(doc *goquery.Document) error {
//	if len(doc.Text()) == 10001 {
//		return errNoPage
//	}
//	return nil
//}
