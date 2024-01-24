package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"example.com/service/service/internal/cached"
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
	cached.GlobalCacheManager.Cache.Set("id", string(body), -1)
}

func OrderView(w http.ResponseWriter, r *http.Request) {
	fmt.Println(12312312)
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
