package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	database "example.com/service/service/internal/database"
	_ "github.com/lib/pq"
)

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//getPeople(w, r)
	case http.MethodPost:
		postOrder(w, r)
	case http.MethodPatch:
		//updatePeople(w, r)
	case http.MethodDelete:
		//deletePeople(w, r)
	default:
		http.Error(w, "Invalid http method", http.StatusMethodNotAllowed)
	}
}

func postOrder(w http.ResponseWriter, r *http.Request) {
	// Чтение тела POST-запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Закрытие тела запроса, чтобы избежать утечек памяти
	defer r.Body.Close()

	// Преобразование данных в строку
	text := strings.Split(string(body), "=")

	// Вывод текста в консоль (или куда-либо еще, в зависимости от вашей логики)
	fmt.Println("Received text:", text[1])

	db, err := database.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.GetDatabase(text[1]); err != nil {
		log.Fatal(err)
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Чтение HTML из файла
	htmlContent, err := ioutil.ReadFile("G:\\Стажировка\\service\\internal\\html\\find.html")
	if err != nil {
		// Если произошла ошибка при чтении файла, возвращаем HTTP 500
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type как text/html
	w.Header().Set("Content-Type", "text/html")

	// Пишем HTML-код в тело ответа
	w.Write(htmlContent)
}
