package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kcpetersen111/iris/server/persist"
)

type PlatformRoutes struct {
	DB *persist.DBInterface
}

type PlatformUserId struct {
	Id string `json:"platformID"`
}

func (p PlatformRoutes) GetPlatform(w http.ResponseWriter, r *http.Request) {
	// var plat persist.Platform
	// log.Println("succ")
	var plat PlatformUserId
	err := json.NewDecoder(r.Body).Decode(&plat)
	if err != nil {
		err := fmt.Errorf("Error in reading body: %v", err)
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	platforms, err := p.DB.GetPlatform(plat.Id)
	if err != nil {
		err := fmt.Errorf("Error in getting data from the database: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(platforms); err != nil {
		log.Printf("%v\n", err)
	}
}

type CreatePlatform struct {
	Name         string   `json:"platformName"`
	UserId       []string `json:"UserId"`
	CreatorEmail string   `json:"email"`
}

func (p *PlatformRoutes) CreatePlatform(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting to create new platform")
	var plat CreatePlatform
	err := json.NewDecoder(r.Body).Decode(&plat)
	if err != nil {
		err := fmt.Errorf("Error in reading body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	err = p.DB.CreatePlatform(plat.Name, append(plat.UserId, plat.CreatorEmail))
	if err == fmt.Errorf("User does not exist") {
		err := fmt.Errorf("User does not exist: %v", err)
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
		return
	}

	if err != nil {
		err := fmt.Errorf("Error in getting data from the database: %v", err)
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	// for _, id := range plat.UserId {

	// 	err = p.DB.CreatePlatform(plat.Name, id)
	// 	if err != nil {
	// 		err := fmt.Errorf("Error in getting data from the database: %v", err)
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		json.NewEncoder(w).Encode(err)
	// 		return
	// 	}
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Println("successful creation of platform")
}
