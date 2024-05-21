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
	"time"
)

type FilmStorage interface {
	//CRUD
	GetFilms(w http.ResponseWriter, r *http.Request)
	GetFilmById(w http.ResponseWriter, r *http.Request)
	AddFilms(w http.ResponseWriter, r *http.Request)
	UpdateFilms(w http.ResponseWriter, r *http.Request)
	UpdateLikes(w http.ResponseWriter, r *http.Request)
	DeleteLikes(w http.ResponseWriter, r *http.Request)
	GetPopularFilms(w http.ResponseWriter, r *http.Request)
}

type InMemoryFilmStorage struct {
	models.Film
}

var Films = make(map[int]InMemoryFilmStorage)

func (f InMemoryFilmStorage) GetFilms(w http.ResponseWriter, r *http.Request) {
	var values []InMemoryFilmStorage
	for _, v := range Films {
		values = append(values, v)
	}
	log.Println("values = ", values)
	json.NewEncoder(w).Encode(values)
}

func (f InMemoryFilmStorage) GetFilmById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmId, err := strconv.Atoi(vars["film_id"])
	fmt.Println("GetFilmById | filmId(1)", filmId)

	if err != nil {
		fmt.Println("GetFilmById | filmId(2)", filmId)
		http.Error(w, "GetFilmById | Invalid film ID", http.StatusBadRequest)
		return
	}
	log.Println("GetFilmById|Get film by Id = ", Films[filmId])
	json.NewEncoder(w).Encode(Films[filmId])
}

func (f InMemoryFilmStorage) AddFilms(w http.ResponseWriter, r *http.Request) {
	film := InMemoryFilmStorage{}
	json.NewDecoder(r.Body).Decode(&film)

	if validateStructure(&film, w) != nil {
		return
	}
	if validateReleaseDate(&film, w) != nil {
		return
	}

	film.ID = IncreaseCounterFilmId()
	Films[film.ID] = film
	json.NewEncoder(w).Encode(film)
}

func (f InMemoryFilmStorage) UpdateFilms(w http.ResponseWriter, r *http.Request) {
	film := InMemoryFilmStorage{}
	json.NewDecoder(r.Body).Decode(&film)

	if validateStructure(&film, w) != nil {
		return
	}

	if validateReleaseDate(&film, w) != nil {
		return
	}

	if CheckFilmByID(film, w) != nil {
		return
	}

	Films[film.ID] = film
	json.NewEncoder(w).Encode(film)
}

func (f InMemoryFilmStorage) UpdateLikes(w http.ResponseWriter, r *http.Request) {
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

	film1 := Films[filmId]
	user1 := Users[userId]

	if CheckFilmByID(film1, w) != nil || CheckUserByID(user1, w) != nil {
		return
	}

	if !utils.Contains(film1.Likes, userId) {
		film1.Likes = append(film1.Likes, userId)
		Films[filmId] = film1

	}

	json.NewEncoder(w).Encode(film1)
}

func (f InMemoryFilmStorage) DeleteLikes(w http.ResponseWriter, r *http.Request) {
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
	user1 := Users[userId]
	film1 := Films[filmId]

	if CheckUserByID(user1, w) != nil || CheckFilmByID(film1, w) != nil {
		return
	}

	if (film1.Likes != nil) && len(film1.Likes) > 0 {
		updatedLikes := utils.Remove(film1.Likes, filmId)
		film1.Likes = updatedLikes
		Films[userId] = film1
	}

	json.NewEncoder(w).Encode(film1)
}

func (f InMemoryFilmStorage) GetPopularFilms(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	likesCountParam := queryParams.Get("count")
	_, err := strconv.Atoi(likesCountParam)
	if err != nil {
		http.Error(w, "Invalid count", http.StatusBadRequest)
		return
	}
}

func IncreaseCounterFilmId() int {
	return len(Films) + 1
}

func validateReleaseDate(film *InMemoryFilmStorage, w http.ResponseWriter) error {
	minDate := time.Date(1895, 12, 28, 0, 0, 0, 0, time.UTC)

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

func validateStructure(film *InMemoryFilmStorage, w http.ResponseWriter) error {
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

func CheckFilmByID(film InMemoryFilmStorage, w http.ResponseWriter) error {
	val, ok := Films[film.ID]
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
