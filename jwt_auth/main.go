package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Id       int     `json:"_id"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Token struct {
	Token string `json:"token"`
}

var (
	client *mongo.Client
)

var userList []User

func main() {
	godotenv.Load()
	connect()
	router := mux.NewRouter()
	router.HandleFunc("/api", getHealth).Methods("GET")
	router.HandleFunc("/api/register", registerUser).Methods("POST")
	router.HandleFunc("/api/login", loginUser).Methods("POST")
	router.HandleFunc("/api/users", getUsers).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbclient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_DB_CONNECTION")))

	if err != nil {
		panic("could not connect")
	}

	client = dbclient
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(bodyBytes, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if newUser.Email == nil {
		createBadRequestResponse(w, http.StatusBadRequest, "Email is required")
	}
	if newUser.Password == nil {
		createBadRequestResponse(w, http.StatusBadRequest, "Password is required")
	}

	passwordInput := *newUser.Password
	newUser.Password = hashAndSalt([]byte(passwordInput))
	userList = append(userList, newUser)

	if userExistsInDb(*newUser.Email) {
		createBadRequestResponse(w, http.StatusBadRequest, "User already exists")
	}

	id := insertUserIntoDb(newUser)

	fmt.Println(id)

	w.WriteHeader(http.StatusCreated)
	return
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fetchedUser := getUserFromDb(*user.Email)

	err = bcrypt.CompareHashAndPassword([]byte(*fetchedUser.Password), []byte(*user.Password))

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := CreateToken(*fetchedUser.Email)

	json.NewEncoder(w).Encode(Token{token})
	w.WriteHeader(http.StatusOK)
	return
}

func createBadRequestResponse(w http.ResponseWriter, errorCode int, message string) {
	w.Header().Set("Content-Type", "applicatoin/json")
	w.WriteHeader(errorCode)

	var response ErrorResponse
	response.Error = http.StatusText(errorCode)
	response.Message = message

	json.NewEncoder(w).Encode(response)
	return
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userList)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func hashAndSalt(password []byte) *string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		panic("error")
	}
	stringHash := string(hash)
	return &stringHash
}

func getUserFromDb(email string) User {
	collection := client.Database("auth").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var fetchedUser User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&fetchedUser)

	if err != nil {
		panic("User could not be inserted into Db")
	}

	return fetchedUser
}

func insertUserIntoDb(newUser User) interface{} {
	collection := client.Database("auth").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.M{"email": newUser.Email, "password": newUser.Password})

	if err != nil {
		panic("User could not be inserted into Db")
	}

	id := res.InsertedID
	return id
}

func userExistsInDb(email string) bool {
	collection := client.Database("auth").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"email": email})

	if err != nil {
		panic("Error checking if user exists in db")
	}

	return count > 0
}

func CreateToken(email string) string {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["email"] = email
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return ""
	}

	return token
}
