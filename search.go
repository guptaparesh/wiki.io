package searcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

    "github.com/guptaparesh/wiki.io/searcher/internal"
	"github.com/guptaparesh/wiki.io/searcher/kvpstore"

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
    defer myutil.LogElapsed(myutil.TrackTime("RunQuery"))
    
	if searchTerm == "" {
		fmt.Println("Empty search term so skipping query execution")
		return nil
	}

    var err error
    var body []byte

    body = kvpstore.GetFromCache(searchTerm)
    if body == nil {
        query := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&list=search&srsearch=%s&srprop=title&format=json", url.QueryEscape(searchTerm))
        body, err = ExecuteGetRequest(query)
        handleError(err)
        kvpstore.Put(searchTerm, body)
    } else {
        log.Println("Found", searchTerm, "in cache")
    }
    
	var wikiRes SrchQueryResponse
	if err := json.Unmarshal(body, &wikiRes); err != nil {
		panic(err)
	}

	// srs := make([]string, 0)
	// for _, sr := range wikiRes.Query.Search {
	// 	srs = append(srs, sr.Title)
	// }
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
