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

	//------MUX
	dt := time.Now()
	fmt.Println("db", dt)
	router := mux.NewRouter()

	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/users/{userId}", controllers.GetUserById).Methods("GET")
	router.HandleFunc("/users", controllers.AddUsers).Methods("POST")
	router.HandleFunc("/users", controllers.UpdateUsers).Methods("PUT")

	router.HandleFunc("/films", controllers.GetFilms).Methods("GET")
	router.HandleFunc("/film/{filmId}", controllers.GetFilmsByID).Methods("GET")
	router.HandleFunc("/films", controllers.AddFilms).Methods("POST")
	router.HandleFunc("/films", controllers.UpdateFilms).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))

	//mux := http.NewServeMux()
	//mux.HandleFunc(`/users/`, getUsers)
	//mux.HandleFunc(`/users/`, putUsers)
	//mux.HandleFunc(`/users/`, postUsers)

	//http.ListenAndServe(":8080", mux)
	//func gin() {
	//	------GIN
	//	router := gin.Default()
	//	router.GET("/films", getFilms)
	//	router.GET("/films/:id", getFilmID)
	//	router.POST("/films", postFilms)
	//
	//	router.GET("/users/", getUsers)
	//	router.POST("/users/", postUsers)
	//	router.PUT("/users/", putUsers)
	//
	//	router.Run("localhost:8080")
	//	}
}

//func getFilms(c *gin.Context) {
//	c.IndentedJSON(http.StatusOK, films)
//}
//
//func getFilmID(c *gin.Context) {
//	//Эта строка извлекает значение параметра "id" из запроса.
//	//В Gin параметры маршрута (route parameters) могут быть доступны через метод Param() объекта gin.Context.
//	id := c.Param("id")
//
//	//Этот цикл перебирает каждый элемент в срезе films, где каждый элемент представляет собой отдельный альбом.
//	for _, a := range films {
//		//В этой части кода мы сравниваем идентификатор альбома a.ID с идентификатором, полученным из параметра запроса id.
//		if a.ID == id {
//			//Если найден фильм с идентификатором, совпадающим с запрошенным,
//			//мы возвращаем информацию об этом фильме в формате JSON с кодом статуса HTTP 200 OK.
//			c.IndentedJSON(http.StatusOK, a)
//			return
//		}
//	}
//	//Если ни один альбом не совпал с запрошенным идентификатором, то ошибка 404 и сообщение
//	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "film not found"})
//}
//
//func postFilms(c *gin.Context) {
//	var newFilm Film
//
//	// Вызывает BindJSON, чтобы привязать полученный JSON к newFilm
//	if err := c.BindJSON(&newFilm); err != nil {
//		return
//	}
//	//Добавляет новый фильм к слайсу
//	films = append(films, newFilm)
//	c.IndentedJSON(http.StatusCreated, newFilm)
//}
//
//func getUsers(c *gin.Context) {
//	c.IndentedJSON(http.StatusOK, users)
//}
//
//func postUsers(c *gin.Context) {
//	var newUser User
//	if err := c.BindJSON(&newUser); err != nil {
//		return
//	}
//	users = append(users, newUser)
//	c.IndentedJSON(http.StatusCreated, newUser)
//}
//
//func putUsers(c *gin.Context) {
//	var updUser User
//	if err := c.BindJSON(&updUser); err != nil {
//		c.AbortWithError(http.StatusBadRequest, err)
//		return
//	}
//}

//func getUsers(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Метод не GET", http.StatusMethodNotAllowed)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(users)
//}

//func putUsers(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPut {
//		http.Error(w, "Метод не PUT", http.StatusMethodNotAllowed)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(users)
//}
//
//func postUsers(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Метод не POST", http.StatusMethodNotAllowed)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(users)
//}
