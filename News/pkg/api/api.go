package api

import (
	"GoNews/pkg/storage"
	"github.com/gorilla/mux"
	"net/http"
)

// Программный интерфейс сервера GoNews
type API struct {
	db     storage.Interface
	router *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.router.HandleFunc("/news", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/{id}", api.DetailPost).Methods(http.MethodGet, http.MethodOptions)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}
