package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go-filmorate/model"
	"go-filmorate/types"
	"log"
	"net/http"
	"strings"
	"time"
)

var users = make(map[int]model.User)

// GetUsers GET users /users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Gets all users")
	values := []model.User{}
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
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}

	log.Println("Birthday after JSON decoding:", user.Birthday)

	validate := validator.New()
	validate.RegisterValidation("validLogin", validLogin)
	//validate.RegisterValidation("validBirthday", validBirthday)

	if user.Name == "" {
		user.Name = user.Login
	}

	if user.Birthday.After(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		responseError := types.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Birthday cannot be in the future"
		json.NewEncoder(w).Encode(responseError)
		return
	}

	err = validate.Struct(user)
	if err = validate.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseError := types.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Bad Request"
		json.NewEncoder(w).Encode(responseError)
		return
	} else {
		user.ID = IncreaseCounterUserId()
	}

	users[user.ID] = user

	log.Println("Add user = ", user)
	json.NewEncoder(w).Encode(user)
}

// UpdateUsers PUT users /users/{userId} (/users/1) + request body
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	log.Println("Update user by id = ", user.ID)
	json.NewDecoder(r.Body).Decode(&user)

	val, ok := users[user.ID]
	if !ok {
		log.Println("nil elem = ", val)
		w.WriteHeader(http.StatusNotFound)
		responseError := types.ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "User Not Found"
		json.NewEncoder(w).Encode(responseError)
		return
	}

	users[user.ID] = user
	json.NewEncoder(w).Encode(user)
}

func IncreaseCounterUserId() int {
	return len(users) + 1
}

func validLogin(fl validator.FieldLevel) bool {
	fmt.Println("Validating login...")

	value := fl.Field().String()
	value2 := " "
	if strings.Contains(value, value2) {
		return false
	}
	return true
}
