package storagies

import (
	"encoding/json"
	"errors"
	"fmt"
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

type InMemoryUserStorage struct {
	models.User
}

var Users = make(map[int]InMemoryUserStorage)

func (u InMemoryUserStorage) GetUsers(w http.ResponseWriter, r *http.Request) {
	var userList []InMemoryUserStorage
	for _, user := range Users {
		userList = append(userList, user)
	}
	json.NewEncoder(w).Encode(userList)
}

func (u InMemoryUserStorage) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(Users[userID])
}

func (u InMemoryUserStorage) AddUsers(w http.ResponseWriter, r *http.Request) {
	var user = InMemoryUserStorage{}

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
	Users[user.ID] = user

	json.NewEncoder(w).Encode(user)
}

func (u InMemoryUserStorage) UpdateUsers(w http.ResponseWriter, r *http.Request) {
	var user = InMemoryUserStorage{}

	json.NewDecoder(r.Body).Decode(&user)

	if err := CheckUserByID(user, w); err != nil {
		return
	}

	Users[user.ID] = user
	json.NewEncoder(w).Encode(user)
}

func (u InMemoryUserStorage) GetFriendsByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	user1 := Users[userId]

	if CheckUserByID(user1, w) != nil {
		return
	}

	userFriendIds := user1.Friends
	userFriend := []InMemoryUserStorage{}
	for _, elem := range userFriendIds {
		userFriend = append(userFriend, Users[elem])
	}
	log.Println("GetFriendByUserId | userFriend = ", userFriend)
	json.NewEncoder(w).Encode(userFriend)
}

func (u InMemoryUserStorage) GetCommonFriendId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])

	if err != nil {
		return
	}
	userOther, err := strconv.Atoi(vars["friend_id"])

	if err != nil {
		return
	}

	user1 := Users[userId]
	user2 := Users[userOther]

	if CheckUserByID(user1, w) != nil || CheckUserByID(user2, w) != nil {
		return
	}

	common := utils.CommonId(user1.Friends, user2.Friends)
	userCommonFriend := []InMemoryUserStorage{}
	for _, elem := range common {
		userCommonFriend = append(userCommonFriend, Users[elem])
	}

	json.NewEncoder(w).Encode(userCommonFriend)
}

func (u InMemoryUserStorage) UpdateFriends(w http.ResponseWriter, r *http.Request) {
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
	fmt.Printf("userId", &userId)
	fmt.Printf("friendId", friendId)

	user1 := Users[userId]
	user2 := Users[friendId]

	fmt.Printf("user1", &user1)
	fmt.Printf("user2", &user2)

	if CheckUserByID(user1, w) != nil || CheckUserByID(user2, w) != nil {
		return
	}

	if !utils.Contains(user1.Friends, friendId) {
		user1.Friends = append(user1.Friends, friendId)
		Users[userId] = user1
	}

	if !utils.Contains(user2.Friends, userId) {
		user2.Friends = append(user2.Friends, userId)
		Users[friendId] = user2
	}

	json.NewEncoder(w).Encode(user2)
}

func (u InMemoryUserStorage) DeleteFriendById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	friendId, err := strconv.Atoi(vars["friend_id"])
	if err != nil {
		return
	}

	user1 := Users[userId]
	user2 := Users[friendId]

	if CheckUserByID(user1, w) != nil || CheckUserByID(user2, w) != nil {
		return
	}

	if (user1.Friends != nil) && len(user1.Friends) > 0 {
		updatedFriends := utils.Remove(user1.Friends, friendId)
		user1.Friends = updatedFriends
		Users[userId] = user1
	}

	if (user2.Friends != nil) && len(user2.Friends) > 0 {
		updatedFriends := utils.Remove(user2.Friends, userId)
		user2.Friends = updatedFriends
		Users[friendId] = user2
	}

	json.NewEncoder(w).Encode(user1)
}

func IncrementUserIDCounter() int {
	return len(Users) + 1
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

func validateBirthday(user *InMemoryUserStorage, w http.ResponseWriter) error {
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

func CheckUserByID(user InMemoryUserStorage, w http.ResponseWriter) error {
	val, ok := Users[user.ID]
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
