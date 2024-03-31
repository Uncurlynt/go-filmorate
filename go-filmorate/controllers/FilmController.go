package controllers

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Film struct {
	ID          int        `json:"id"`
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required,lte=200"`
	ReleaseDate CustomTime `json:"releaseDate" validate:"required"`
	Duration    int        `json:"duration" validate:"required,gt=0"`
}

var films = make(map[int]Film)

// GET films
func GetFilms(w http.ResponseWriter, r *http.Request) {
	log.Println("Gets all films")
	values := []Film{}
	for _, v := range films {
		values = append(values, v)
	}
	log.Println("values = ", values)
	json.NewEncoder(w).Encode(values)
}

// GET films by ID
func GetFilmsByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmId := vars["filmId"]
	log.Println("Get film by id = ", filmId)
	json.NewEncoder(w).Encode(filmId)
}

// POST films
func AddFilms(w http.ResponseWriter, r *http.Request) {

	film := Film{}
	json.NewDecoder(r.Body).Decode(&film)
	film.ID = IncreaseCounterFilmId()
	validate := validator.New()

	log.Println("ReleaseDate after JSON decoding:", film.ReleaseDate)

	minDate := time.Date(1895, 12, 28, 0, 0, 0, 0, time.UTC)
	log.Println("minDate :", minDate)

	if film.ReleaseDate.Before(minDate) {
		w.WriteHeader(http.StatusBadRequest)
		responseError := ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Release date must be after 28 December 1895"
		json.NewEncoder(w).Encode(responseError)
		return
	}

	err := validate.Struct(film)
	if err = validate.Struct(film); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//Лог тела функции
		responseError := ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "User Not Found"
		json.NewEncoder(w).Encode(responseError)
		return
		//} else {
		//	film.ID = IncreaseCounterUserId()
	}
	films[film.ID] = film
	log.Println("Add film - ", film.ID)
	json.NewEncoder(w).Encode(film)
}

// PUT films
func UpdateFilms(w http.ResponseWriter, r *http.Request) {
	film := Film{}
	log.Println("Update film by id = ", film.ID)
	json.NewDecoder(r.Body).Decode(&film)

	validate := validator.New()

	err := validate.Struct(film)
	if err = validate.Struct(film); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//Лог тела функции
		responseError := ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "User Not Found"
		json.NewEncoder(w).Encode(responseError)
		return
	}
	val, ok := films[film.ID]
	if !ok {
		log.Println("nil elem = ", val)
		w.WriteHeader(http.StatusNotFound)
		responseError := ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "Film Not Found"
		json.NewEncoder(w).Encode(responseError)
		return
	}

	films[film.ID] = film
	json.NewEncoder(w).Encode(film)
}

func IncreaseCounterFilmId() int {
	return len(films) + 1
}
