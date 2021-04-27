package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

var api = "https://api"
var hotCache = cache.New(30*time.Second, 10*time.Minute)
var coldCache = cache.New(-1, -1)

func loadData(w http.ResponseWriter, req *http.Request) {

	var url = api + req.URL.Path + string('?') + req.URL.RawQuery

	cachedResponse, found := hotCache.Get(url)

	if found {
		fmt.Println("Cached: " + url)
		coldCache.SetDefault(url, cachedResponse)
		fmt.Fprintf(w, cachedResponse.(string))
	} else {

		client := http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get(url)

		if err != nil {
			cachedResponse, found = coldCache.Get(url)

			if found {
				fmt.Fprintf(w, cachedResponse.(string))
			} else {
				fmt.Fprintf(w, "") // todo : send response with error code
			}

			return
		}

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			hotCache.Set(url, bodyString, cache.DefaultExpiration)
			coldCache.SetDefault(url, bodyString)
			fmt.Println("HTTP Status is in the 2xx range")
		} else {
			cachedResponse, found = coldCache.Get(url) // handle case when cold cache not found
			fmt.Println("Error HTTP Status code")
			bodyString = cachedResponse.(string)
		}

		fmt.Fprint(w, bodyString)
	}

}

func main() {
	http.HandleFunc("/", loadData)
	http.ListenAndServe(":8811", nil)
}
