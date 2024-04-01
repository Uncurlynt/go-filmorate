package controllers

import (
	"encoding/json"
	"go-filmorate/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var users = make(map[int]models.User)
var checkUser = true

// GetUsers GET users /users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Gets all users")
	values := []models.User{}
	for _, v := range users {
		values = append(values, v)
	}
	log.Println("values = ", values)
	json.NewEncoder(w).Encode(values)
}

// GetUserById GET users /users/{userId} (/users/1)
func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	log.Println("Get user by id = ", userId)
	json.NewEncoder(w).Encode(userId)
}

// AddUsers POST users /users + request body
func AddUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	log.Println("Birthday after JSON decoding:", user.Birthday)

	validate := validator.New()
	validate.RegisterValidation("validLogin", validLogin)

	//validateName(&user)
	if user.Name == "" {
		user.Name = user.Login
	}

	if !validateBirthday(&user, w) {
		return
	}

	if !validateEmail(&user, w) {
		return
	}

	user.ID = IncreaseCounterUserId()
	users[user.ID] = user
	log.Println("Add user = ", user)
	json.NewEncoder(w).Encode(user)
}

// UpdateUsers PUT users /users/{userId} (/users/1) + request body
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	log.Println("Update user by id = ", user.ID)
	json.NewDecoder(r.Body).Decode(&user)

	if !checkUserByID(&user, w) {
		return
	}

	users[user.ID] = user
	json.NewEncoder(w).Encode(user)
}

func IncreaseCounterUserId() int {
	return len(users) + 1
}

func validLogin(fl validator.FieldLevel) bool {
	log.Println("Validating login...")

	value := fl.Field().String()
	value2 := " "
	if strings.Contains(value, value2) {
		checkUser = false
		return false
	}
	checkUser = true
	return true
}

func validateBirthday(user *models.User, w http.ResponseWriter) bool {

	if user.Birthday.After(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Birthday cannot be in the future"
		json.NewEncoder(w).Encode(responseError)
		checkUser = false
		return false
	}
	return true
}

func validateEmail(user *models.User, w http.ResponseWriter) bool {
	validate := validator.New()
	err := validate.Struct(user)
	if err = validate.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Bad Request"
		json.NewEncoder(w).Encode(responseError)
		//checkUser = false
		return false
	}
	return true
}

func checkUserByID(user *models.User, w http.ResponseWriter) bool {
	val, ok := users[user.ID]
	if !ok {
		log.Println("nil elem = ", val)
		w.WriteHeader(http.StatusNotFound)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "User Not Found"
		json.NewEncoder(w).Encode(responseError)
		return false
	}
	return true
}
