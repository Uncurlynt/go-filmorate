package storagies

import (
	"net/http"
)

// FilmStorage Интерфейс надо заполнить функциями из контроллеров (Имя и тип возвращаемого значения)
type FilmStorage interface {
	GetFilms() error
	GetFilmById() error
	AddFilms() error
	UpdateFilm() error
	UpdateLikes() error
	DeleteLikes() error
	GetPopularFilms() error
}

// InMemoryFilmStorage Определение структуры
type InMemoryFilmStorage struct {
	//Здесь добавить поля из FilmModel
}

// Get Метод Get для структуры InMemoryFilmStorage (не понимаю, что здесь делать)
func (us InMemoryFilmStorage) Update(w http.ResponseWriter, r *http.Request) {
	//	Здесь пишем логиек функций
}

//func NewInMemoryFilmStorage() *InMemoryFilmStorage {
//}

type UserStorage interface {
}

type InMemoryUserStorage struct {
}
