package main

import (
	"net/http"
	"time"
	"log"
	"io"
	"fmt"
)

// var TargetURL = "https://4pda.ru"
var TargetURL = "https://yandex.ru"

func sendResp ( w http.ResponseWriter, url string) {
	fmt.Printf("url %s \n", url)

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal( err )
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)


	if err != nil {
		log.Fatal( err )
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func getResp ( w http.ResponseWriter, r *http.Request) {
	fmt.Printf("r.URL %s \n", r.URL)

	// if r.URL.Path == "/"  {
	// 	urlstr := fmt.Sprintf("%s", TargetURL)
	// 	sendResp ( w, urlstr)

	// } else {
	// 	urlstr := fmt.Sprintf("/%s", r.URL)
	// 	sendResp ( w, urlstr)
	// }

	urlstr := fmt.Sprintf("%s?%s", TargetURL, r.URL.Path)
	sendResp ( w, urlstr)

	// fmt.Printf( " urlstr %s \n", urlstr)
	// fmt.Printf("r.URL.RawQuery %s \n", r.URL.RawQuery)
	// fmt.Printf("r.URL.Path %s \n", r.URL.Path)


}

func initServer() {

    s := &http.Server{
        Addr:           ":8080",
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

	// if err != nil {
	// 	log.Fatal(err)
	// }

	http.HandleFunc("/", getResp)
	s.ListenAndServe()

}

func main () {

    initServer()
}
