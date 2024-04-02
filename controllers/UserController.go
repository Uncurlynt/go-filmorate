package controllers

import (
	"encoding/json"
	"errors"
	"go-filmorate/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var users = make(map[int]models.User)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Gets all users")
	values := []models.User{}
	for _, v := range users {
		values = append(values, v)
	}
	log.Println("values = ", values)
	json.NewEncoder(w).Encode(values)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		panic(err)
	}
	log.Println("Get user by id = ", users[userId])
	json.NewEncoder(w).Encode(users[userId])
}

func AddUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	log.Println("Birthday after JSON decoding:", user.Birthday)

	validate := validator.New()
	validate.RegisterValidation("validLogin", validLogin)

	if user.Name == "" {
		user.Name = user.Login
	}

	if validateBirthday(&user, w) != nil {
		return
	}

	if validateEmail(&user, w) != nil {
		return
	}

	user.ID = IncreaseCounterUserId()
	users[user.ID] = user
	log.Println("Add user = ", user)
	json.NewEncoder(w).Encode(user)
}

func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	log.Println("Update user by id = ", user.ID)
	json.NewDecoder(r.Body).Decode(&user)

	if checkUserByID(&user, w) != nil {
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
		return false
	}
	return true
}

func validateBirthday(user *models.User, w http.ResponseWriter) error {
	if user.Birthday.After(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Birthday cannot be in the future"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("Birthday cannot be in the future")
	}
	return nil
}

func validateEmail(user *models.User, w http.ResponseWriter) error {
	validate := validator.New()
	err := validate.Struct(user)
	if err = validate.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Bad Request"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("Bad Request")
	}
	return nil
}

func checkUserByID(user *models.User, w http.ResponseWriter) error {
	val, ok := users[user.ID]
	if !ok {
		log.Println("nil elem = ", val)
		w.WriteHeader(http.StatusNotFound)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "User Not Found"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("User Not Found")
	}
	return nil
}
