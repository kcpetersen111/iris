package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/kcpetersen111/iris/server/persist"
	"github.com/kcpetersen111/iris/server/persist/users"

	"golang.org/x/crypto/bcrypt"
)

type UserRoutes struct {
	DB *persist.DBInterface
}

func CreateUserRouter(db *persist.DBInterface) mux.Router {
	userStruct := UserRoutes{
		DB: db,
	}
	router := mux.NewRouter()

	router.HandleFunc("/user", userStruct.SignUp).Methods("POST")

	return *router
}

// func (u *UserRoutes) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	var input users.MessageJson
// 	rawRequest, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		log.Printf("Error reading request: %v", err)
// 	}
// 	err = json.Unmarshal(rawRequest, &input)
// 	if err != nil {
// 		log.Printf("Error unmarshaling request: %v", err)
// 	}

// 	users.CreateUser(input.Username, input.Password)

// 	log.Printf("dbtest request body %+v,", input)

// 	w.WriteHeader(http.StatusOK)
// 	// Write([]byte(http.StatusOK))

// }

//user auth guide

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

func (u UserRoutes) SignUp(w http.ResponseWriter, r *http.Request) {
	// connection := GetDatabase()
	// defer Closedatabase(connection)

	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// var err error
		err := fmt.Errorf("Error in reading body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	// var dbuser User
	// connection.Where("email = ?", user.Email).First(&dbuser)
	dbuser, err := user.GetUserByEmail(u.DB)
	if err != nil {
		err := fmt.Errorf("Error in DB: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	//checks if email is already register or not
	if dbuser.Email != "" {
		// var err error
		err = fmt.Errorf("Email already in use")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	user.Password, err = GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	//insert user details in database
	err = user.CreateUser(u.DB)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateJWT(email, role string) (string, error) {
	// secretkey := os.Getenv("jwt")
	var mySigningKey = []byte("JWTTOKEN")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func (u UserRoutes) SignIn(w http.ResponseWriter, r *http.Request) {

	var authdetails users.Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		err = fmt.Errorf("Error in reading body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	authuser, err := authdetails.GetUserByEmail(u.DB)
	// fmt.Println(authdetails)
	// connection.Where("email = ?", authdetails.Email).First(&authuser)

	if authuser.Email == "" {
		err = fmt.Errorf("Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	check := CheckPasswordHash(authdetails.Password, authuser.Password)

	// fmt.Printf("HASH: %v, P1: %v, P2: %v \n", check, authdetails.Password, authuser.Password)
	if !check {
		err = fmt.Errorf("Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := GenerateJWT(authuser.Email, authuser.Role)
	if err != nil {
		err = fmt.Errorf("Failed to generate token: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var token Token
	token.Email = authuser.Email
	token.Role = authuser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// type middleware interface {
// 	Middleware(handler http.Handler) http.Handler
// }

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(r.Header["Token"])

		if r.Header["Token"] == nil {
			err := fmt.Errorf("No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte("JWTTOKEN")

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err error
			err = fmt.Errorf("Your Token has been expired: %v", err)
			json.NewEncoder(w).Encode(err)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {

				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {

				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return
			}
		}

		err = fmt.Errorf("Not Authorized")
		json.NewEncoder(w).Encode(err)
	}
}
