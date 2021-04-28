package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"handlers"

	"github.com/patrickmn/go-cache"
)

var api = os.Getenv("CS_API")
var hotCache = cache.New(45*time.Second, 5*time.Minute)
var coldCache = cache.New(-1, -1) // cache that does not expire , only update

func loadData(w http.ResponseWriter, req *http.Request) {

	var url = api + req.URL.Path + string('?') + req.URL.RawQuery

	cachedResponse, found := hotCache.Get(url)

	if found {
		fmt.Printf("Page from HOT cache: %s\n", url)
		coldCache.SetDefault(url, cachedResponse.(string))
		fmt.Fprint(w, cachedResponse.(string))
		return
	} else {

		client := http.Client{
			Timeout: 60 * time.Second,
		}

		resp, err := client.Get(url)

		if err != nil {

			cachedResponse, found = coldCache.Get(url)

			if found {
				fmt.Printf("Page from COLD cache: %s\n", url)
				fmt.Fprint(w, cachedResponse.(string))
			} else {
				fmt.Println("Worst case, nothing sent to browser")
				fmt.Fprint(w, "") // todo : send response with error code
			}

			return
		}

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			hotCache.Set(url, bodyString, cache.DefaultExpiration)
			coldCache.SetDefault(url, bodyString)
			fmt.Printf("Requested page, set coldCache and hotCache %s\n", url)
		} else {
			fmt.Printf("Page from COLD cache #2: %s\n", url)
			cachedResponse, found = coldCache.Get(url) // handle case when cache not found
			bodyString = cachedResponse.(string)
		}

		fmt.Fprint(w, bodyString)
	}

}

func main() {
	http.HandleFunc("/", loadData)
	http.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "pong")
	})
	http.HandleFunc("/health-check", handlers.HealthCheckHandler)
	log.Fatal(http.ListenAndServe(":8811", nil))
}
