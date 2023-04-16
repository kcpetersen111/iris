package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kcpetersen111/iris/server/persist"
	"github.com/kcpetersen111/iris/server/persist/users"

	"golang.org/x/crypto/bcrypt"
)

type UserRoutes struct {
	DB *persist.DBInterface
}

// func CreateUserRouter(db *persist.DBInterface) mux.Router {
// 	userStruct := UserRoutes{
// 		DB: db,
// 	}
// 	router := mux.NewRouter()

// 	router.HandleFunc("/user", userStruct.SignUp).Methods("POST")

// 	return *router
// }

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
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(err)
		return
	}

	user.Password, err = GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	//insert user details in database
	uuid, err := user.CreateUser(u.DB)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	user.UserID = uuid.String()
	jwt, err := GenerateJWT(user.Email, "user")
	user.Name = jwt
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
	claims["exp"] = time.Now().Add(time.Minute * 300).Unix()

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

	var token users.User
	token.Email = authuser.Email
	token.Role = authuser.Role
	token.Name = validToken
	token.UserID = authuser.UserID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		log.Printf("%v\n", err)
	}

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
		// log.Printf("User authenticating")
		// fmt.Println(r.Header)

		if r.Header["Authorization"] == nil {
			err := fmt.Errorf("No Token Found")
			log.Println(err)
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte("JWTTOKEN")

		token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			// var err error
			err = fmt.Errorf("Your Token has been expired: %v", err)
			log.Println(err)
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
		log.Println(err)
		json.NewEncoder(w).Encode(err)
	}
}
