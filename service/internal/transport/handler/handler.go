package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	cached "example.com/service/service/internal/cached"
	database "example.com/service/service/internal/database"
	_ "github.com/lib/pq"
)

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		postOrder(w, r)
	case http.MethodPatch:
	case http.MethodDelete:
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

	db, err := database.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//order, err := db.GetDatabase(string(body))
	//if err != nil {
	//	log.Fatal(err)
	//}

	fmt.Println(cached.GlobalCacheManager.Cache.Get(string(body)))
}

func OrderView(w http.ResponseWriter, r *http.Request) {
	// Чтение HTML из файла
	htmlContent, err := ioutil.ReadFile("G:\\Стажировка\\service\\internal\\html\\order.html")
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
