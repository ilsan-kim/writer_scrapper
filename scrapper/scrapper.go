package scrapper

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type novelDataStruct struct {
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
var novelDataList []novelDataStruct

// Scrape naver webtoon ID by num
func Scrape(num int) {
	c := make(chan novelDataStruct)
	excel := excelize.NewFile()
	sheet := excel.NewSheet("Sheet1")
	excel.SetActiveSheet(sheet)

	defer excel.SaveAs("NovelData.xlsx")
	excel.SetCellValue("Sheet1", "A1", "URL")
	excel.SetCellValue("Sheet1", "B1", "제목")
	excel.SetCellValue("Sheet1", "C1", "최초 작성일")
	excel.SetCellValue("Sheet1", "D1", "설명")
	excel.SetCellValue("Sheet1", "E1", "장르")
	excel.SetCellValue("Sheet1", "F1", "작가 필명")
	excel.SetCellValue("Sheet1", "G1", "작가 이메일")
	success := 2
	for i:=num ; i<num+2000 ; i++ {
		queryURL := baseURL + strconv.Itoa(i)
		go InsertDataToStruct(queryURL, c)
		time.Sleep(time.Microsecond * 3500)
	}
	for i:=num ; i<num+2000 ; i++{
		novelData := <-c

		if novelData.url != "" {
			novelDataList = append(novelDataList, novelData)
			urlCell := fmt.Sprintf("A%s", strconv.Itoa(success))
			titleCell := fmt.Sprintf("B%s", strconv.Itoa(success))
			pubDateCell := fmt.Sprintf("C%s", strconv.Itoa(success))
			descCell := fmt.Sprintf("D%s", strconv.Itoa(success))
			genreCell := fmt.Sprintf("E%s", strconv.Itoa(success))
			nicknameCell := fmt.Sprintf("F%s", strconv.Itoa(success))
			emailCell := fmt.Sprintf("G%s", strconv.Itoa(success))
			excel.SetCellValue("Sheet1", urlCell, novelData.url)
			excel.SetCellValue("Sheet1", titleCell, novelData.title)
			excel.SetCellValue("Sheet1", pubDateCell, novelData.firstPubDate)
			excel.SetCellValue("Sheet1", descCell, novelData.novelDesc)
			excel.SetCellValue("Sheet1", genreCell, novelData.genre)
			excel.SetCellValue("Sheet1", nicknameCell, novelData.writerNickname)
			excel.SetCellValue("Sheet1", emailCell, novelData.writerEmail)
			success = success + 1
		}
	}
	time.Sleep(time.Second * 2)
	fmt.Println(novelDataList)
}

func InsertDataToStruct(url string, c chan<- novelDataStruct) {
	pages := GetPages(url)
	if pages != nil {
		desc, email := GetNovelDesc(pages)
		fmt.Printf("Requesting %s : SUCCESS \n", url)
		c <- novelDataStruct{
			url:            url,
			title:          GetNovelTitle(pages),
			firstPubDate:   GetPubDate(pages),
			novelDesc:      desc,
			genre:          GetGenre(pages),
			writerNickname: GetWriterNickname(pages),
			writerEmail:    email,
		}
	}
	fmt.Printf("Requesting %s : DELETED \n", url)
	c <- novelDataStruct{}
}

func GetPages(url string) *goquery.Document {
	res, err := http.Get(url)
	checkErr(url, err)
	checkStatus(url, res.StatusCode)
	res.Header.Add("Connection", "close")
	defer res.Body.Close()

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