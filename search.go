package searcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

type SrchQueryResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Sroffset int    `json:"sroffset"`
		Continue string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Searchinfo struct {
			Totalhits int `json:"totalhits"`
		} `json:"searchinfo"`
		Search []struct {
			Ns        int       `json:"ns"`
			Title     string    `json:"title"`
			Pageid    int       `json:"pageid"`
			Size      int       `json:"size"`
			Wordcount int       `json:"wordcount"`
			Snippet   string    `json:"snippet"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"search"`
	} `json:"query"`
}

type MostViewedResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Pvimoffset int    `json:"pvimoffset"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Mostviewed []struct {
			Ns    int    `json:"ns"`
			Title string `json:"title"`
			Count int    `json:"count"`
		} `json:"mostviewed"`
	} `json:"query"`
}

func ExecuteGetRequest(query string) ([]byte, error) {
	if query == "" {
		return nil, errors.New("non-empty query is required")
	}

	resp, err := http.Get(query)
	handleError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	return body, err
}

func RunQuery(searchTerm string) *SrchQueryResponse {
	if searchTerm == "" {
		fmt.Println("Empty search term so skipping query execution")
		return nil
	}
	body, err := ExecuteGetRequest(fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&list=search&srsearch=%s&srprop=title&format=json", searchTerm))
	handleError(err)

	var wikiRes SrchQueryResponse
	if err := json.Unmarshal(body, &wikiRes); err != nil {
		panic(err)
	}

	srs := make([]string, 0)
	for _, sr := range wikiRes.Query.Search {
		srs = append(srs, sr.Title)
	}
	return &wikiRes
}

func FindMostViewed(pvimlimit int) map[string]int {
	query := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&list=mostviewed&pvimmetric=pageviews&pvimlimit=%d&format=json", pvimlimit)
	body, err := ExecuteGetRequest(query)
	handleError(err)

	var mvr MostViewedResponse
	if err := json.Unmarshal(body, &mvr); err != nil {
		panic(err)
	}

	titleToViewCountMap := make(map[string]int)
	for _, mv := range mvr.Query.Mostviewed {
		titleToViewCountMap[mv.Title] = mv.Count
	}

	return titleToViewCountMap
}
