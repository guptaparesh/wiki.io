package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"git.mills.io/prologic/bitcask"
	"github.com/gorilla/mux"
	"github.com/guptaparesh/wiki.io/searcher"
	"github.com/guptaparesh/wiki.io/searcher/internal"
	"github.com/guptaparesh/wiki.io/searcher/kvpstore"
)

type HandlerContext struct {
	bcDB *bitcask.Bitcask
}

func NewHandlerContext(db *bitcask.Bitcask) *HandlerContext {	
	if db == nil {
		log.Fatalln("nil Bitcask database")
	}

	return &HandlerContext{db}
}

func (ctx *HandlerContext) TermSearchHandler(rw http.ResponseWriter, req *http.Request) {
	defer myutil.LogElapsed(myutil.TrackTime("TermSearchHandler"))

	vars := mux.Vars(req)
	searchTerm := vars["term"]
	fmt.Println("Search term", searchTerm)
	//rw.Write([]byte("Hello from TermSearchHandler"))

	var err error
    var body []byte

    body = kvpstore.GetIfExists(ctx.bcDB, searchTerm)
    if body == nil {
        query := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&list=search&srsearch=%s&srprop=title&format=json", 
		url.QueryEscape(searchTerm))
        body, err = searcher.ExecuteGetRequest(query)
        handleError(err, "")
        kvpstore.PutInDb(ctx.bcDB, searchTerm, body)
    } else {
        log.Println("Found", searchTerm, "in cache")
    }

	rw.Write(body)
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg,err)
	}
}

func main() {
	addr := ":8080"
	fmt.Println("Starting wiki server")
	DB, err := bitcask.Open("/tmp/db")
	defer DB.Close()

	handleError(err, "Can't open bitcask database")

	//construct the handler context
	handlerCtx := NewHandlerContext(DB)

	router := mux.NewRouter()
	router.HandleFunc("/search/{term}", handlerCtx.TermSearchHandler)

	log.Printf("wiki search server is listening at %s...", addr)
    log.Fatal(http.ListenAndServe(addr, router))
}