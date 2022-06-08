package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type resultsResponse struct{
	Info interface{} `json:"Info"`
	Result struct {
		Exams []Exam `json:"Exams"`
	} `json:"Result"`
}

type Results []Exam

type Exam struct{
	Name string `json:"Subject"`
	Mark int `json:"TestMark"`
	HasResult bool `json:"HasResult"`
}

const apiUrl string = "https://checkege.rustest.ru/api/exam"

func GetResults(participantId string) (Results, error) {
	cookiejar, _ := cookiejar.New(nil)
	cookie := http.Cookie{
		Name:       "Participant",
		Value:      participantId}
	client := http.Client{
		Jar: cookiejar,
	}
	urlObj, _ := url.Parse(apiUrl)
	client.Jar.SetCookies(urlObj, []*http.Cookie{&cookie})

	resp, err := client.Get(apiUrl)
	if err != nil{
		return nil, fmt.Errorf("request to api failed: %s", err)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	resultsResp := &resultsResponse{}
	if err := json.Unmarshal(data, resultsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %s", err)
	}

	return resultsResp.Result.Exams, nil
}