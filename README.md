# GO client and server cmd applications for exploring wikimedia via API 
* They both invoke [wikipedia api](https://www.mediawiki.org/wiki/API:Main_page) and cache results to avoid refetching them from wikipedia
* [bitcask](https://git.mills.io/prologic/bitcask) is used as a locally persistenty cache
## To use the client 
```shell
$ go install github.com/guptaparesh/wiki.io/searcher/cmd/wiki-explorer/
$ wiki-explorer
```

## To use the HTTP server 
```shell
$ go install github.com/guptaparesh/wiki.io/searcher/cmd/wiki-server/
$ wiki-server
```
