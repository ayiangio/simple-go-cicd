package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Item is struct for export
type Item struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Rating float32 `json:"rating"`
	Stock  int     `json:"stock"`
	Price  float32 `json:"price"`
}

var (
	items  = []Item{}
	nextID = 1
	mu     sync.Mutex
)

func main() {
	http.HandleFunc("/items", loggingMiddleware(itemsHandler))
	http.HandleFunc("/items/", loggingMiddleware(itemHandler))
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logResponseWriter := &responseLogger{ResponseWriter: w}
		next(logResponseWriter, r)
		fmt.Printf("%s\n", logResponseWriter.body)
	}
}

type responseLogger struct {
	http.ResponseWriter
	statusCode int
	body       string
}

func (rl *responseLogger) WriteHeader(statusCode int) {
	rl.statusCode = statusCode
	rl.ResponseWriter.WriteHeader(statusCode)
}

func (rl *responseLogger) Write(b []byte) (int, error) {
	rl.body = string(b)
	return rl.ResponseWriter.Write(b)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItems(w, r)
	case http.MethodPost:
		createItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		deleteItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// JSONResponse is a struct used for sending JSON responses.
type JSONResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	ClientIP   string      `json:"ip,omitempty"`
	Method     string      `json:"method,omitempty"`
	Path       string      `json:"path,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func jsonResponse(w http.ResponseWriter, statusCode int, message string, data interface{}, ip string, method string, path string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResponse := JSONResponse{
		StatusCode: statusCode,
		Message:    message,
		ClientIP:   ip,
		Method:     method,
		Path:       path,
		Data:       data,
	}
	json.NewEncoder(w).Encode(jsonResponse)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()
	jsonResponse(w, http.StatusOK, "Success", items, r.RemoteAddr, r.Method, r.URL.Path)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	item.ID = nextID
	nextID++
	items = append(items, item)
	w.Header().Set("Content-Type", "application/json")
	jsonResponse(w, http.StatusCreated, "Item created successfully", item, r.RemoteAddr, r.Method, r.URL.Path)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")
	var found bool
	mu.Lock()
	defer mu.Unlock()
	for i, item := range items {
		if strconv.Itoa(item.ID) == id {
			items = append(items[:i], items[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		jsonResponse(w, http.StatusNotFound, "Item not found", nil, r.RemoteAddr, r.Method, r.URL.Path)
		return
	}
	jsonResponse(w, http.StatusNoContent, "Item deleted successfully", nil, r.RemoteAddr, r.Method, r.URL.Path)
}
