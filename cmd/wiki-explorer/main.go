package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"git.mills.io/prologic/bitcask"
	"github.com/guptaparesh/wiki.io/searcher"
	"github.com/guptaparesh/wiki.io/searcher/kvpstore"
)

func main() {
	var err error
	// acCmd 	= autocomplete page title search
	// sCmd		= wiki pages search
	// vcCmd 	= page view count search

	// setup Bitcask
	kvpstore.DB, err = bitcask.Open("/tmp/db")
	if err != nil {
		log.Fatal(err)
	}
	defer kvpstore.DB.Close()

	sCmd := flag.NewFlagSet("search", flag.ExitOnError)
	sTerm := sCmd.String("t", "", "wikipedia search term")

	vcCmd := flag.NewFlagSet("vc", flag.ExitOnError)
	vcLimit := vcCmd.Int("limit", 10, "how many top page view counts to return")

	if len(os.Args) < 2 {
		fmt.Println("Expected either 'search' or 'vc' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "search":
		sCmd.Parse(os.Args[2:])
		runSearch(*sTerm)
	case "vc":
		vcCmd.Parse(os.Args[2:])
		runVcCmd(*vcLimit)
	default:
		fmt.Println("Expected either 'search' or 'vc' subcommands")
		os.Exit(1)
	}

}

func runSearch(term string) {
	//defer LogElapsed(TrackTime("runSearch"))
	fmt.Println("Searching for", term)
	searchResponse := searcher.RunQuery(term)
	if searchResponse != nil {
		fmt.Println("TotalHits: ", searchResponse.Query.Searchinfo.Totalhits)

		for _, s := range searchResponse.Query.Search {
			fmt.Println("Title:", s.Title)
		}
	}
}

func runVcCmd(vcLimit int) {
	titleToViewsCountMap := searcher.FindMostViewed(vcLimit)

	titles := make([]string, 0, len(titleToViewsCountMap))
	for t := range titleToViewsCountMap {
		titles = append(titles, t)
	}
	sort.Strings(titles)

	//viewCount
	sort.Slice(titles, func(i, j int) bool {
		return titleToViewsCountMap[titles[i]] > titleToViewsCountMap[titles[j]]
	})

	fmt.Println("Title                     ViewCount")
	fmt.Println("------------------------------------")
	for _, title := range titles {
		fmt.Printf("%-20s  %10d\n", title, titleToViewsCountMap[title])
	}
}
