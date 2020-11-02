package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"
)

var targetURL = os.Getenv("URL")

type Example struct {
	Greet string `json:"greet"`
}

//var greet []Example

func LogInformation(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req := fmt.Sprintf("Request to %s with %s ", r.URL, r.Method)
		log.Println(req)
		h.ServeHTTP(w, r)

	})

}

func main() {

	//greet = append(greet, Example{Greet: "Hello Universe!"})

	//FileServer := http.FileServer(http.Dir("."))
	//http.Handle("/", FileServer)
	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1", hellouniverse).Methods("GET")
	mux.Use(LogInformation)

	http.ListenAndServe(":9091", mux)
}

func hellouniverse(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(targetURL)
	if err != nil {
		//fmt.Printf("Error!", err)
		errmsg := fmt.Sprintf("Error while http.Get: %s", err.Error())
		log.Println(errmsg)
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	log.Println("Server status: ", resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Printf("Error!", err)
		errmsg := fmt.Sprintf("Error while read body: %s", err.Error())
		log.Println(errmsg)
		http.Error(w, errmsg, http.StatusServiceUnavailable)
		return

	}

	responseString := string(body)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	} else {

		Hello := Example{Greet: fmt.Sprintf("Greeting to World: %s", responseString)}

		w.Header().Set("Content-Type", "applicaton/json")
		json.NewEncoder(w).Encode(Hello)

	}
}
