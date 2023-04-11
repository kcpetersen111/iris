package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kcpetersen111/iris/server/persist"
	"github.com/kcpetersen111/iris/server/ports/websockets"
	"github.com/rs/cors"
)

type Server interface {
	Serve() (Server, error)
}
type IrisServer struct {
	Address string
	DB      *persist.DBInterface
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

func (s *IrisServer) dbtest(w http.ResponseWriter, r *http.Request) {
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

	s.DB.InsertTest(input.Message)
	// fmt.Println(r.Header["Role"])
	w.Write([]byte("message: " + input.Message))
}

func (s *IrisServer) Serve() {

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/test", IsAuthorized(s.dbtest)).Methods("POST")

	// User methods
	// ur := router.PathPrefix("/user").Subrouter()
	// userRouter := CreateUserRouter(s.DB)
	user := UserRoutes{
		DB: s.DB,
	}

	router.HandleFunc("/createUser", user.SignUp).Methods("POST")
	router.HandleFunc("/user", user.SignIn).Methods("POST")

	hub := websockets.NewHub(s.DB)

	router.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		websockets.StartCall(hub, w, r)
	})

	platform := PlatformRoutes{
		DB: s.DB,
	}
	router.HandleFunc("/getPlatform", IsAuthorized(platform.GetPlatform)).Methods("POST")
	router.HandleFunc("/platform", IsAuthorized(platform.CreatePlatform)).Methods("POST")

	message := MessageRoutes{
		DB: s.DB,
	}

	router.HandleFunc("/getMessages", IsAuthorized(message.GetMessage)).Methods("POST")
	router.HandleFunc("/message", IsAuthorized(message.PostMessage)).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	srv := &http.Server{
		Handler:      c.Handler(router),
		Addr:         s.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func NewIrisServer(address string, db *persist.DBInterface) *IrisServer {
	return &IrisServer{
		Address: address,
		DB:      db,
	}
}
