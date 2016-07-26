//package main
//
//import (
//	"github.com/nats-io/nats"
//	"fmt"
//	"time"
//	"os"
//	"log"
//	"io/ioutil"
//	"encoding/json"
//	"net/http"
//)
//
//func main (){
//	nc, err := nats.Connect(nats.DefaultURL)
//	if err != nil {
//		panic(fmt.Sprintf("Cannot connect to NATS: %s", err))
//	}
//	defer nc.Close()
//
//	resp, err := nc.Request("dataservice.get.1234.lastscan", []byte{}, time.Second)
//	fmt.Printf("Body recieved: %s", string(resp.Data))
//}

package main

import (
"net/http"
"fmt"
"log"
"io/ioutil"
"nextevolution/capitulum/types"
	dataService "nextevolution/data-service/types"
"encoding/json"
"os"
	"github.com/nats-io/nats"
	"time"
)

var config types.Config

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to NATS: %s", err))
	}
	defer nc.Close()

	if len(os.Args) <= 1  || os.Args[1] == ""{
		log.Panic("Please supply a config file path like: ./mock config.json")
	}
	configPath := os.Args[1]

	rawConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to read config file %s", configPath))
	}

	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Panic(fmt.Sprintf("unable to unmarshal config file %s", configPath))
	}

	http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		GiveResponseFile("responses/GET_alive.txt",w)
	})

	// List Sales
	http.HandleFunc("/api/sales", func(w http.ResponseWriter, r *http.Request){
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "405 method not allowed")
			return
		}

		if r.Header.Get("AuthToken") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401 unauthorized")
			return
		}

		if r.Header.Get("AuthToken") == "bad_token" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401 unauthorized")
			return
		}

		resp, err := nc.Request("dataservice.get.1234.lastscan", []byte{}, time.Second)
		if err != nil {
			log.Fatalf("Got NATS error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 internal server error")
			return
		}

		sas := dataService.SellerAlbumScan{}
		err = json.Unmarshal(resp.Data, &sas)
		if err != nil {
			log.Fatalf("Unable to unmarshal response body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 internal server error")
			return
		}

		salesList := types.ConvertDataSales(sas)
		salesListJson, err := json.Marshal(salesList)
		if err != nil {
			log.Fatalf("Unable to marshal sales list: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 internal server error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(salesListJson))
	})

	// Facebook Login
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request){
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "405 method not allowed")
			return
		}

		//unmarshal data
		var loginReq types.FbLoginReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "400 bad request - Can't read body")
			return
		}

		err = json.Unmarshal(body, &loginReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "400 bad request - Can't read json")
			return
		}

		if loginReq.FbToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "400 bad request - empty token")
			return
		}

		if loginReq.FbToken == "bad_token" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401 unauthorized")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		GiveResponseFile("responses/POST_api_login.json",w)
	})

	http.Handle("/", http.FileServer(http.Dir(config.StaticFilePath)))

	log.Printf("Listening on port: %d", config.Port)
	log.Printf("Static file path: %s", config.StaticFilePath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), EveryCall(http.DefaultServeMux)))
}

func GiveResponseFile(filename string, w http.ResponseWriter){
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("reading sample file (%s)failed", filename))
	}
	fmt.Fprintf(w, string(file))
}

func EveryCall(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		// determine origin header
		origin := r.Header.Get("Origin")

		for _, orgn := range config.Origins {
			log.Printf("Request Origin: %s Trusted Origin: %s",origin, orgn)
			if orgn == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Response-Type, If-Modified-Since")
		w.Header().Set("Access-Control-Expose-Header", "Content-disposition, X-Pagination-Current-Page, X-Pagination-Page-Count, X-Pagination-Per-Page, X-Pagination-Total-Count, Link")
		if r.Method == "OPTIONS" {
			//return just headers
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func getSales (){

}