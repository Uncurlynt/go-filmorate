package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

type User struct {
	ID       int        `json:"id"`
	Email    string     `json:"email" validate:"required,email"`
	Login    string     `json:"login" validate:"required,validLogin"`
	Name     string     `json:"name"`
	Birthday CustomTime `json:"birthday"`
}

type CustomTime struct {
	time.Time
}

type ResponseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var users = make(map[int]User)

// GetUsers GET users /users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Gets all users")
	values := []User{}
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
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// обработка ошибки декодирования JSON
		return
	}

	log.Println("Birthday after JSON decoding:", user.Birthday) // Добавляем логирование

	validate := validator.New()
	validate.RegisterValidation("validLogin", validLogin)
	//validate.RegisterValidation("validBirthday", validBirthday)

	if user.Name == "" {
		user.Name = user.Login
	}

	if user.Birthday.After(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		responseError := ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Birthday cannot be in the future"
		json.NewEncoder(w).Encode(responseError)
		return
	}

	err = validate.Struct(user)
	if err = validate.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseError := ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Bad Request"
		json.NewEncoder(w).Encode(responseError)
		return
	} else {
		user.ID = IncreaseCounterUserId()
	}

	users[user.ID] = user //По кючу user.ID вставляем значения user

	log.Println("Add user = ", user)
	json.NewEncoder(w).Encode(user)
}

// UpdateUsers PUT users /users/{userId} (/users/1) + request body
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	user := User{}
	log.Println("Update user by id = ", user.ID)
	json.NewDecoder(r.Body).Decode(&user)

	val, ok := users[user.ID]
	if !ok {
		log.Println("nil elem = ", val)
		w.WriteHeader(http.StatusNotFound)
		responseError := ResponseError{}
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

func (t CustomTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	} else {
		return []byte(`"` + t.Format("2006-01-02") + `"`), nil
	}
}

func (t *CustomTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	log.Println("Unmarshalled time string:", timeStr) // Добавляем логирование
	parsedTime, err := time.Parse("2006-01-02", timeStr)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

func validLogin(fl validator.FieldLevel) bool {
	fmt.Println("Validating login...") // Добавляем логирование

	value := fl.Field().String()
	value2 := " "
	if strings.Contains(value, value2) {
		return false
	}
	return true
}
