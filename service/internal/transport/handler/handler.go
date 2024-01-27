package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

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
	cached.GlobalCacheManager.Set("id", string(body), -1)
}

func OrderView(w http.ResponseWriter, r *http.Request) {
	// Чтение HTML из файла
	//htmlContent, err := ioutil.ReadFile("G:\\Стажировка\\service\\internal\\html\\order.html")
	htmlContent, err := ioutil.ReadFile("internal\\html\\order.html")
	if err != nil {
		// Если произошла ошибка при чтении файла, возвращаем HTTP 500
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type как text/html
	w.Header().Set("Content-Type", "text/html")

	// Parse the HTML template
	tmpl, err := template.New("order").Parse(string(htmlContent))
	if err != nil {
		panic(err)
	}

	flag := true

	if cachedValueId, foundId := cached.GlobalCacheManager.Get("id"); foundId {
		cached.GlobalCacheManager.Delete("id")
		id := fmt.Sprintf("%v", cachedValueId)

		if cachedValueOrder, foundOrder := cached.GlobalCacheManager.Get(id); foundOrder {
			var outputBuffer bytes.Buffer

			err = tmpl.Execute(&outputBuffer, cachedValueOrder)
			if err != nil {
				panic(err)
			}

			w.Write(outputBuffer.Bytes())
		} else {
			flag = false
		}
	} else {
		flag = false
	}

	if !flag {
		// Чтение HTML из файла
		//htmlContent, err := ioutil.ReadFile("G:\\Стажировка\\service\\internal\\html\\error.html")
		htmlContent, err := ioutil.ReadFile("internal\\html\\error.html")
		if err != nil {
			// Если произошла ошибка при чтении файла, возвращаем HTTP 500
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовок Content-Type как text/html
		w.Header().Set("Content-Type", "text/html")

		w.Write(htmlContent)
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Чтение HTML из файла
	//htmlContent, err := ioutil.ReadFile("G:\\Стажировка\\service\\internal\\html\\find.html")
	htmlContent, err := ioutil.ReadFile("internal\\html\\find.html")
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
