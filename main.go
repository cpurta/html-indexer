package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	influxClient "github.com/influxdata/influxdb/client/v2"
	"gopkg.in/redis.v2"
	"golang.org/x/net/html"
)

var (
	UrlQueueHost string
	UrlQueuePort string

	IndexResultHost string
	IndexResultPort string

	IndexResultUser string
	IndexResultPass string

	InfluxHost string
	InfluxPort string
)

func main() {
	log.Println("Loading environment variables...")
	loadEnvironmentVariables()

	// url queue client
	urlClient, err := ConnectToRedis(UrlQueueHost, UrlQueuePort)

	// result mysql client
	indexClient, err := ConnectToMysql(IndexResultUser, IndexResultPass, IndexResultHost, IndexResultPort)

	// influx database client
	influx, err := influxClient.NewHTTPClient(influxClient.HTTPConfig{
		Addr:     fmt.Sprintf("%s:%s", InfluxHost, InfluxPort),
		Username: "root",
		Password: "temppwd",
		Timeout:  time.Second * 10,
	})

	if err != nil {
		log.Fatalln("Error connecting to our resources:", err.Error())
	}

	// get all keys
	keys := urlClient.Keys("*").Val()

	var wg sync.WaitGroup
	for i := 0; i < len(keys); i++ {
		wg.Add(1)
		go func(urls *redis.Client, indexes *redis.Client, dbKey string) {
			defer wg.Done()

			IndexURLs(urls, indexes, dbKey)
		}(urlClient, indexClient, keys[i])
	}

	wg.Wait()
	log.Println("Finished indexing all urls in Redis queues")
}

func IndexURLs(urls *redis.Client, db *sql.DB, dbKey string) {
	for url := urls.RPop(dbKey); url != nil {
		keyCount = make(map[string]int)

		// grab the response from the site...
		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		if resp.Body != nil {
			defer resp.Body.Close()

			tokenizer := html.NewTokenizer(resp.Body)

			for {
				tt := tokenizer.Next()
				switch tt {
				case html.ErrorToken:
					return
				case html.TextToken:
					text := string(tokenizer.Text())

					// sanitize our data... mostly remove special characters
					sanitized := strings.Replace(text, "!@#$%^&*()-+={}[]\\|/?<>,.~`", "", -1)

					pieces := strings.Split(sanitized, " ")

					// we will need to make sure that these values are
					for p := range pieces {
						// no point in keeping track of emtpy strings from the split
						if p != "" {
							if keyCount[p] == nil {
								keyCount[p] = 1
							} else {
								keyCount[p]++
							}
						}
					}
				case html.StartTagToken:
					return
				} // end switch
			} // end for

			writeResults(keyCount, url, db)
		} // end if
	}// end for urls
}

func writeResults(results map[string]int, url string, db *sql.DB) error {
	return nil
}

func loadEnvironmentVariables() {
	UrlQueueHost = os.GetEnv("URLQUEUE_PORT_6379_TCP_ADDR")
	UrlQueuePort = os.GetEnv("URLQUEUE_PORT_6379_TCP_PORT")
	IndexResultHost = os.GetEnv("WEBINDEXDB_PORT_3308_TCP_ADDR")
	IndexResultPort = os.GetEnv("WEBINDEXDB_PORT_3308_TCP_PORT")
	IndexResultUser = os.Getenv("WEBINDEXDB_USERNAME")
	IndexResultPass = os.Getenv("WEBINDEXDB_PASSWORD")
	InfluxHost = os.Getenv("INFLUXDB_PORT_8086_TCP_ADDR")
	InfluxPort = os.Getenv("INFLUXDB_PORT_8086_TCP_PORT")
}

func ConnectToRedis(host, port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Printf("Connecting to redis server %s:%s: %s %s\n", host, port, pong, err.Error())

	return client
}

func ConnectToMysql(user, pass, host, port string) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, IndexResultHost, IndexResultPort))
}
