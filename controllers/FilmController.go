package controllers

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go-filmorate/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

var films = make(map[int]models.Film)

func GetFilms(w http.ResponseWriter, r *http.Request) {
	log.Println("Gets all films")
	values := []models.Film{}
	for _, v := range films {
		values = append(values, v)
	}
	log.Println("values = ", values)
	json.NewEncoder(w).Encode(values)
}

func GetFilmById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmId, err := strconv.Atoi(vars["filmId"])
	if err != nil {
		panic(err)
	}
	log.Println("Get film by Id = ", films[filmId])
	json.NewEncoder(w).Encode(films[filmId])
}

func AddFilms(w http.ResponseWriter, r *http.Request) {
	film := models.Film{}
	json.NewDecoder(r.Body).Decode(&film)

	if validateStructure(&film, w) != nil {
		return
	}
	if validateReleaseDate(&film, w) != nil {
		return
	}

	film.ID = IncreaseCounterFilmId()
	films[film.ID] = film
	json.NewEncoder(w).Encode(film)
}

func UpdateFilms(w http.ResponseWriter, r *http.Request) {
	film := models.Film{}
	log.Println("Update film by id = ", film.ID)
	json.NewDecoder(r.Body).Decode(&film)

	if validateStructure(&film, w) != nil {
		return
	}

	if validateReleaseDate(&film, w) != nil {
		return
	}

	if checkFilmById(&film, w) != nil {
		return
	}

	films[film.ID] = film
	json.NewEncoder(w).Encode(film)

}

func IncreaseCounterFilmId() int {
	return len(films) + 1
}

func validateReleaseDate(film *models.Film, w http.ResponseWriter) error {
	minDate := time.Date(1895, 12, 28, 0, 0, 0, 0, time.UTC)
	log.Println("minDate :", minDate)
	log.Println("ReleaseDate after JSON decoding:", film.ReleaseDate)

	if film.ReleaseDate.Before(minDate) {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Release date must be after 28 December 1895"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("Release date must be after 28 December 1895")
	}
	return nil
}

func validateStructure(film *models.Film, w http.ResponseWriter) error {
	validate := validator.New()
	err := validate.Struct(film)
	if err = validate.Struct(film); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Validation Error"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("Validation Error")
	}
	return nil
}

func checkFilmById(film *models.Film, w http.ResponseWriter) error {
	val, ok := films[film.ID]
	if !ok {
		log.Println("nil elem = ", val)
		w.WriteHeader(http.StatusNotFound)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "Film Not Found"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("Film Not Found")
	}
	return nil
}
