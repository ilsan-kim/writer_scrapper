package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var baseURL = "https://novel.naver.com/challenge/list?novelId="

func main() {
	for i:=0 ; i<100 ; i++{
		queryURL := baseURL + strconv.Itoa(i)
		fmt.Println(GetStatus(queryURL))
	}
}

func GetStatus(url string) int {
	res, err := http.Get(url)
	checkErr(url, err)
	checkStatus(url, res.StatusCode)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	return res.StatusCode
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