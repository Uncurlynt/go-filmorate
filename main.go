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

	dt := time.Now()
	fmt.Println("db", dt)
	router := mux.NewRouter()

	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/users/{userId}", controllers.GetUserById).Methods("GET")
	router.HandleFunc("/users", controllers.AddUsers).Methods("POST")
	router.HandleFunc("/users", controllers.UpdateUsers).Methods("PUT")

	router.HandleFunc("/films", controllers.GetFilms).Methods("GET")
	router.HandleFunc("/film/{filmId}", controllers.GetFilmById).Methods("GET")
	router.HandleFunc("/films", controllers.AddFilms).Methods("POST")
	router.HandleFunc("/films", controllers.UpdateFilms).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
