package controllers

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
	"sort"
	"strconv"
	"time"
)

type Film struct {
	ID          int              `json:"id"`
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description" validate:"required,lte=200"`
	ReleaseDate utils.CustomTime `json:"releaseDate" validate:"required"`
	Duration    int              `json:"duration" validate:"required,gt=0"`
	Likes       []int            `json:"likes"`
}

var films = make(map[int]Film)

type FilmStorage interface {
	GetFilms(w http.ResponseWriter, r *http.Request)
	GetFilmById(w http.ResponseWriter, r *http.Request)
	AddFilms(w http.ResponseWriter, r *http.Request)
	UpdateFilms(w http.ResponseWriter, r *http.Request)
	UpdateLikes(w http.ResponseWriter, r *http.Request)
	DeleteLikes(w http.ResponseWriter, r *http.Request)
	GetPopularFilms(w http.ResponseWriter, r *http.Request)
}

func (f *Film) GetFilms(w http.ResponseWriter, r *http.Request) {
	var values []Film
	for _, v := range films {
		values = append(values, v)
	}
	log.Println("values = ", values)
	json.NewEncoder(w).Encode(values)
}

func (f *Film) GetFilmById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmId, err := strconv.Atoi(vars["film_id"])
	fmt.Println("GetFilmById | filmId(1)", filmId)

	if err != nil {
		fmt.Println("GetFilmById | filmId(2)", filmId)
		http.Error(w, "GetFilmById | Invalid film ID", http.StatusBadRequest)
		return
	}
	log.Println("GetFilmById|Get film by Id = ", films[filmId])
	json.NewEncoder(w).Encode(films[filmId])
}

func (f *Film) AddFilms(w http.ResponseWriter, r *http.Request) {
	film := Film{}
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

func (f *Film) UpdateFilms(w http.ResponseWriter, r *http.Request) {
	film := Film{}
	//log.Println("Update film by id = ", film.ID)
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

func (f *Film) UpdateLikes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmId, err := strconv.Atoi(vars["film_id"])
	if err != nil {
		http.Error(w, "UpdateLikes | Invalid film ID", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	film1 := films[filmId]
	user1 := users[userId]

	if checkFilmById(&film1, w) != nil || checkUserByID(user1, w) != nil {
		return
	}

	if !Contains(film1.Likes, userId) {
		film1.Likes = append(film1.Likes, userId)
		films[filmId] = film1

	}

	json.NewEncoder(w).Encode(film1)
}

func (f *Film) DeleteLikes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "1Invalid user ID", http.StatusBadRequest)
		return
	}
	filmId, err := strconv.Atoi(vars["film_id"])
	if err != nil {
		http.Error(w, "DeleteLikes | Invalid film ID", http.StatusBadRequest)
		return
	}

	user1 := users[userId]
	film1 := films[filmId]

	if checkUserByID(user1, w) != nil || checkFilmById(&film1, w) != nil {
		return
	}

	if (film1.Likes != nil) && len(film1.Likes) > 0 {
		updatedLikes := Remove(film1.Likes, filmId)
		film1.Likes = updatedLikes
		films[userId] = film1
	}

	json.NewEncoder(w).Encode(film1)
}

func (f *Film) GetPopularFilms(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	likesCountParam := queryParams.Get("count")
	likesCount, err := strconv.Atoi(likesCountParam)
	if err != nil {
		http.Error(w, "Invalid count", http.StatusBadRequest)
		return
	}

	fmt.Println("likesCount:", likesCount)

	filmsSlice := make([]Film, 0, len(films))
	for _, film := range films {
		filmsSlice = append(filmsSlice, film)
	}

	sort.Slice(filmsSlice, func(i, j int) bool {
		film1 := filmsSlice[i]
		film2 := filmsSlice[j]
		return len(film1.Likes) > len(film2.Likes)
	})

	if likesCount > len(filmsSlice) {
		likesCount = len(filmsSlice)
	}

	json.NewEncoder(w).Encode(filmsSlice[:likesCount])
}

func IncreaseCounterFilmId() int {
	return len(films) + 1
}

func validateReleaseDate(film *Film, w http.ResponseWriter) error {
	minDate := time.Date(1895, 12, 28, 0, 0, 0, 0, time.UTC)
	//log.Println("minDate :", minDate)
	//log.Println("ReleaseDate after JSON decoding:", film.ReleaseDate)

	if film.ReleaseDate.Before(minDate) {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Release date must be after 28 December 1895"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("release date must be after 28 December 1895")
	}
	return nil
}

func validateStructure(film *Film, w http.ResponseWriter) error {
	validate := validator.New()
	err := validate.Struct(film)
	if err = validate.Struct(film); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusBadRequest
		responseError.Message = "Validation Error"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("validation Error")
	}
	return nil
}

func checkFilmById(film *Film, w http.ResponseWriter) error {
	val, ok := films[film.ID]
	if !ok {
		log.Println("nil elem = ", val)
		w.WriteHeader(http.StatusNotFound)
		responseError := models.ResponseError{}
		responseError.Status = http.StatusNotFound
		responseError.Message = "Film Not Found"
		json.NewEncoder(w).Encode(responseError)
		return errors.New("film Not Found")
	}
	return nil
}
