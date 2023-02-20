package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kcpetersen111/iris/server/persist"

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

type User struct {
	UserID   string `json:UserID`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

// func GetDatabase() *gorm.DB {
// 	databasename := "userdb"
// 	database := "postgres"
// 	databasepassword := "1312"
// 	databaseurl := "postgres://postgres:" + databasepassword + "@localhost/" + databasename + "?sslmode=disable"
// 	connection, err := gorm.Open(database, databaseurl)
// 	if err != nil {
// 		log.Fatalln("wrong database url")
// 	}

// 	sqldb := connection.DB()

// 	err = sqldb.Ping()
// 	if err != nil {
// 		log.Fatal("database connected")
// 	}

// 	fmt.Println("connected to database")
// 	return connection
// }

// func InitialMigration() {
// 	connection := GetDatabase()
// 	defer Closedatabase(connection)
// 	connection.AutoMigrate(User{})
// }

// func Closedatabase(connection *gorm.DB) {
// 	sqldb := connection.DB()
// 	sqldb.Close()

// }

func (u UserRoutes) SignUp(w http.ResponseWriter, r *http.Request) {
	// connection := GetDatabase()
	// defer Closedatabase(connection)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// var err error
		err := fmt.Errorf("Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	// var dbuser User
	// connection.Where("email = ?", user.Email).First(&dbuser)
	dbuser, err := user.VerifyEmail(u.DB)
	if err != nil {
		err := fmt.Errorf("Error in DB")
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
	connection.Create(&user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(secretkey)
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

func SignIn(w http.ResponseWriter, r *http.Request) {
	connection := GetDatabase()
	defer Closedatabase(connection)

	var authdetails Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		var err Error
		err = SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var authuser User
	connection.Where("email = ?", authdetails.Email).First(&authuser)
	if authuser.Email == "" {
		var err Error
		err = SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	check := CheckPasswordHash(authdetails.Password, authuser.Password)

	if !check {
		var err Error
		err = SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := GenerateJWT(authuser.Email, authuser.Role)
	if err != nil {
		var err Error
		err = SetError(err, "Failed to generate token")
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

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			var err Error
			err = SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte(secretkey)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err Error
			err = SetError(err, "Your Token has been expired")
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
		var reserr Error
		reserr = SetError(reserr, "Not Authorized")
		json.NewEncoder(w).Encode(err)
	}
}
