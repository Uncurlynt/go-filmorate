package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-filmorate/storagies"
	"log"
	"net/http"
	"time"
)

func main() {
	var NewInMemoryUserStorage = storagies.InMemoryUserStorage{}
	var NewInMemoryFilmStorage = storagies.InMemoryFilmStorage{}

	dt := time.Now()
	fmt.Println("db", dt)
	router := mux.NewRouter()

	router.HandleFunc("/users", NewInMemoryUserStorage.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", NewInMemoryUserStorage.GetUserById).Methods("GET")
	router.HandleFunc("/users", NewInMemoryUserStorage.AddUsers).Methods("POST")
	router.HandleFunc("/users", NewInMemoryUserStorage.UpdateUsers).Methods("PUT")
	router.HandleFunc("/users/{id}/friends", NewInMemoryUserStorage.GetFriendsByUserId).Methods("GET")
	router.HandleFunc("/users/{id}/friends/{friend_id}", NewInMemoryUserStorage.UpdateFriends).Methods("PUT")
	router.HandleFunc("/users/{id}/friends/{friend_id}", NewInMemoryUserStorage.DeleteFriendById).Methods("DELETE")
	router.HandleFunc("/users/{id}/friends/common/{friend_id}", NewInMemoryUserStorage.GetCommonFriendId).Methods("GET")

	router.HandleFunc("/films", NewInMemoryFilmStorage.GetFilms).Methods("GET")
	router.HandleFunc("/films", NewInMemoryFilmStorage.AddFilms).Methods("POST")
	router.HandleFunc("/films", NewInMemoryFilmStorage.UpdateFilms).Methods("PUT")
	router.HandleFunc("/films/{film_id}/like/{id}", NewInMemoryFilmStorage.UpdateLikes).Methods("PUT")
	router.HandleFunc("/films/{film_id}/like/{id}", NewInMemoryFilmStorage.DeleteLikes).Methods("DELETE")
	router.HandleFunc("/films/popular", NewInMemoryFilmStorage.GetPopularFilms).Methods("GET")
	router.HandleFunc("/films/{film_id}", NewInMemoryFilmStorage.GetFilmById).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
