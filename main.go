package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func loadData(w http.ResponseWriter, req *http.Request) {
	var url = "https://www.jonathanfielding.com" + req.URL.Path

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

	fmt.Fprint(w, bodyString)
}

func main() {
	http.HandleFunc("/", loadData)
	http.ListenAndServe(":8000", nil)
}
