package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server interface {
	Serve() (Server, error)
}
type IrisServer struct {
	address string
}
type InputJson struct {
	Message string `json:"message"`
}

func ping(w http.ResponseWriter, _r *http.Request) {
	fmt.Println("Hello, Iris!")
	message, err := json.Marshal("Hello, Iris!")
	if err != nil {
		log.Printf("Something went wrong when pinging the server\n")
	}
	w.Write(message)
}

func dbtest(w http.ResponseWriter, r *http.Request) {
	var input InputJson
	rawRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request: %v", err)
	}
	err = json.Unmarshal(rawRequest, &input)
	if err != nil {
		log.Printf("Error unmarshaling request: %v", err)
	}

	log.Printf("dbtest request body %+v,", input)

	w.Write([]byte("message: " + input.Message))
}

func (s *IrisServer) Serve() {

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/test", dbtest).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         s.address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
func NewIrisServer(address string) *IrisServer {
	return &IrisServer{
		address: address,
	}
}
