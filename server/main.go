package main

import (
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kcpetersen111/iris/server/persist"
)

func ping(_w http.ResponseWriter, _r *http.Request) {
	fmt.Println("Hello, Iris")

}

func main() {
	persist.DbSetupConnection(true)

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
