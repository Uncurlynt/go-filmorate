package controllers

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go-filmorate/models"
	"go-filmorate/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID       int              `json:"id"`
	Email    string           `json:"email" validate:"required,email"`
	Login    string           `json:"login" validate:"required,isValidLogin"`
	Name     string           `json:"name"`
	Birthday utils.CustomTime `json:"birthday"`
	Friends  []int            `json:"friends"`
}

var users = make(map[int]User)

type UserStorage interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	AddUsers(w http.ResponseWriter, r *http.Request)
	UpdateUsers(w http.ResponseWriter, r *http.Request)
	GetFriendsByUserId(w http.ResponseWriter, r *http.Request)
	GetCommonFriendId(w http.ResponseWriter, r *http.Request)
	UpdateFriends(w http.ResponseWriter, r *http.Request)
	DeleteFriendById(w http.ResponseWriter, r *http.Request)
}

func (u User) GetUsers(w http.ResponseWriter, r *http.Request) {
	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}
	json.NewEncoder(w).Encode(userList)
}

func (u User) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(users[userID])
}

func (u User) AddUsers(w http.ResponseWriter, r *http.Request) {
	var user = User{}

	json.NewDecoder(r.Body).Decode(&user)

	validate := validator.New()
	validate.RegisterValidation("isValidLogin", isValidLogin)

	if user.Name == "" {
		user.Name = user.Login
	}

	if err := validateBirthday(&user, w); err != nil {
		return
	}

	if !strings.Contains(user.Email, "@") || strings.Contains(user.Email, " ") {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid email format",
		}

		json.NewEncoder(w).Encode(responseError)
		errors.New("Invalid email format")
		return
	}

	user.ID = IncrementUserIDCounter()
	user.Friends = []int{}
	users[user.ID] = user

	json.NewEncoder(w).Encode(user)
}

func (u User) UpdateUsers(w http.ResponseWriter, r *http.Request) {
	var user = User{}

	json.NewDecoder(r.Body).Decode(&user)

	if err := checkUserByID(user, w); err != nil {
		return
	}

	users[user.ID] = user
	json.NewEncoder(w).Encode(user)
}

func (u User) GetFriendsByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	user1 := users[userId]

	if checkUserByID(user1, w) != nil {
		return
	}

	userFriendIds := user1.Friends
	userFriend := []User{}
	for _, elem := range userFriendIds {
		userFriend = append(userFriend, users[elem])
	}
	//user2 := users[userFriend]

	log.Println("GetFriendByUserId | userFriend = ", userFriend)
	json.NewEncoder(w).Encode(userFriend)
}

func (u User) GetCommonFriendId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])

	if err != nil {
		panic(err)
	}
	userOther, err := strconv.Atoi(vars["friend_id"])

	if err != nil {
		panic(err)
	}

	user1 := users[userId]
	user2 := users[userOther]

	if checkUserByID(user1, w) != nil || checkUserByID(user2, w) != nil {
		return
	}

	common := commonId(user1.Friends, user2.Friends)
	userCommonFriend := []User{}
	for _, elem := range common {
		userCommonFriend = append(userCommonFriend, users[elem])
	}
	//НАСРАЛ ЖИРНЮЩИЙ КОСТЫЛЬ
	json.NewEncoder(w).Encode(userCommonFriend)
	//json.NewEncoder(w).Encode(common)
}

func (u User) UpdateFriends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	friendId, err := strconv.Atoi(vars["friend_id"])
	if err != nil {
		http.Error(w, "Invalid friend ID", http.StatusBadRequest)
		return
	}

	user1 := users[userId]
	user2 := users[friendId]

	if checkUserByID(user1, w) != nil || checkUserByID(user2, w) != nil {
		return
	}

	if !Contains(user1.Friends, friendId) {
		user1.Friends = append(user1.Friends, friendId)
		users[userId] = user1
	}

	if !Contains(user2.Friends, userId) {
		user2.Friends = append(user2.Friends, userId)
		users[friendId] = user2
	}

	json.NewEncoder(w).Encode(user2)
}

func (u User) DeleteFriendById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	friendId, err := strconv.Atoi(vars["friend_id"])
	if err != nil {
		panic(err)
	}

	user1 := users[userId]
	user2 := users[friendId]

	if checkUserByID(user1, w) != nil || checkUserByID(user2, w) != nil {
		return
	}

	if (user1.Friends != nil) && len(user1.Friends) > 0 {
		updatedFriends := Remove(user1.Friends, friendId)
		//users[user1.ID].Friends = updatedFriends
		user1.Friends = updatedFriends
		users[userId] = user1
	}

	if (user2.Friends != nil) && len(user2.Friends) > 0 {
		updatedFriends := Remove(user2.Friends, userId)
		user2.Friends = updatedFriends
		users[friendId] = user2
	}

	json.NewEncoder(w).Encode(user1)
}

func IncrementUserIDCounter() int {
	return len(users) + 1
}

func isValidLogin(fl validator.FieldLevel) bool {
	log.Println("Validating login...")

	value := fl.Field().String()
	value2 := " "
	if strings.Contains(value, value2) {
		return false
	}
	return true
}

func validateBirthday(user *User, w http.ResponseWriter) error {
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

func checkUserByID(user User, w http.ResponseWriter) error {
	val, ok := users[user.ID]
	if !ok {
		log.Println("nil elem = ", val)
		w.WriteHeader(http.StatusNotFound)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "User Not Found"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("user Not Found")
	}
	return nil
}

func Contains(s []int, val int) bool {
	for _, item := range s {
		if item == val {
			return true
		}
	}
	return false
}

func commonId(user1 []int, user2 []int) []int {
	seen := make([]int, 1000)
	for i := range user1 {
		seen[user1[i]]++
	}

	res := make([]int, 0)
	for i := range user2 {
		if seen[user2[i]] > 0 {
			res = append(res, user2[i])
			seen[user2[i]] = 0
		}
	}
	return res
}

func Remove(slice []int, s int) []int {
	for i, num := range slice {
		if num == s {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}
