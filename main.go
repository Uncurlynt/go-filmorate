package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-filmorate/controllers"
	"log"
	"net/http"
	"time"
)

func main() {
	var InMemoryUserStorage controllers.UserStorage = controllers.User{}
	var InMemoryFilmStorage controllers.FilmStorage = &controllers.Film{}

	dt := time.Now()
	fmt.Println("db", dt)
	router := mux.NewRouter()

	router.HandleFunc("/users", InMemoryUserStorage.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", InMemoryUserStorage.GetUserById).Methods("GET")
	router.HandleFunc("/users", InMemoryUserStorage.AddUsers).Methods("POST")
	router.HandleFunc("/users", InMemoryUserStorage.UpdateUsers).Methods("PUT")
	router.HandleFunc("/users/{id}/friends", InMemoryUserStorage.GetFriendsByUserId).Methods("GET")
	router.HandleFunc("/users/{id}/friends/{friend_id}", InMemoryUserStorage.UpdateFriends).Methods("PUT")
	router.HandleFunc("/users/{id}/friends/{friend_id}", InMemoryUserStorage.DeleteFriendById).Methods("DELETE")
	router.HandleFunc("/users/{id}/friends/common/{friend_id}", InMemoryUserStorage.GetCommonFriendId).Methods("GET")

	router.HandleFunc("/films", InMemoryFilmStorage.GetFilms).Methods("GET")
	router.HandleFunc("/films", InMemoryFilmStorage.AddFilms).Methods("POST")
	router.HandleFunc("/films", InMemoryFilmStorage.UpdateFilms).Methods("PUT")
	router.HandleFunc("/films/{film_id}/like/{id}", InMemoryFilmStorage.UpdateLikes).Methods("PUT")
	router.HandleFunc("/films/{film_id}/like/{id}", InMemoryFilmStorage.DeleteLikes).Methods("DELETE")
	router.HandleFunc("/films/popular", InMemoryFilmStorage.GetPopularFilms).Methods("GET")
	router.HandleFunc("/films/{film_id}", InMemoryFilmStorage.GetFilmById).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
