package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

var pageCache = cache.New(5*time.Minute, 10*time.Minute)

func loadData(w http.ResponseWriter, req *http.Request) {
	var api = "https://api"
	var url = api + req.URL.Path + string('?') + req.URL.RawQuery

	cachedResponse, found := pageCache.Get(url)

	if found {
		fmt.Println("Cached result: " + url)
		fmt.Fprintf(w, cachedResponse.(string))
	} else {
		resp, err := http.Get(url)

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Println("HTTP Status is in the 2xx range")
		} else {
			fmt.Println("Error HTTP Status code")
		}

		bodyString := string(bodyBytes)

		pageCache.Set(url, bodyString, cache.DefaultExpiration)

		fmt.Fprint(w, bodyString)
	}

}

func main() {
	http.HandleFunc("/", loadData)
	http.ListenAndServe(":8811", nil)
}
