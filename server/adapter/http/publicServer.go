package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func ping(_w http.ResponseWriter, _r *http.Request) {
	fmt.Println("Hello, Iris")

}

func startServer() {

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
