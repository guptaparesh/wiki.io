package kvpstore

import (
	"log"
	"git.mills.io/prologic/bitcask"
)

var DB *bitcask.Bitcask;

func SetupDb() {
	db, err := bitcask.Open("/tmp/db")
	if err != nil {
		log.Fatalln("failed to open db connection", err)
	} else {
		log.Println("Connected to db. Currently has ", db.Len(), "items")
	}
}

func ExampleNgw(k string, v string) string {
	SetupDb()
    defer Cleanup()

    DB.Put([]byte(k), []byte(v))
    val, err := DB.Get([]byte(k))
	if err != nil {
		log.Fatalf("error %v\n", err)
	}
    //log.Printf(string(val))
	return string(val)
}

func Example(k string, v string) string {
	db, err := bitcask.Open("/tmp/db")
    defer db.Close()

    db.Put([]byte(k), []byte(v))
    val, err := db.Get([]byte(k))
	if err != nil {
		log.Fatalf("error %v\n", err)
	}
    log.Printf(string(val))
	return string(val)
}

func Cleanup() error {
	return DB.Close()
}

func GetFromCache(key string) []byte {
	k := []byte(key)
	if DB.Has(k) {
		val, err := DB.Get(k)
		if err != nil {
			log.Fatalf("error %v\n", err)
		}
		return val
	} 
	return nil	
}

func Put(key string, val []byte) {
	k := []byte(key)
	err := DB.Put(k, val)
	if err != nil {
		log.Println("Error putting", key, "in cache")
	}
} 
