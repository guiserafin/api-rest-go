package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type id uuid.UUID

type user struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Biography string `json:"biography"`
}

type postResponse struct {
	UUID string `json:"uuid"`
}

type aplication struct {
	data map[id]user
}

type response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func sendJSON(w http.ResponseWriter, resp response, status int) {

	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(resp)

	if err != nil {
		sendJSON(w, response{Message: "Something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err := w.Write(data); err != nil {
		fmt.Println("Erro ao enviar resposta", err)
		return
	}

}

func HttpHandler() http.Handler {

	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(jsonMiddleware)

	db := aplication{
		data: make(map[id]user),
	}

	primeiro, _ := uuid.NewRandom()

	fmt.Println(primeiro)

	db.data[id(primeiro)] = user{
		FirstName: "Guilherme",
		LastName:  "Carvalho",
		Biography: "Golang dev",
	}

	r.Route("/api", func(r chi.Router) {

		r.Post("/users", handlePostUser(&db))
		r.Get("/users", findAll(db))
		r.Get("/users/{id}", handlegetUser(db))
		r.Delete("/users/{id}", handleDeleteUser(db))
		r.Put("/users/{id}", putUser(db))

	})

	return r
}

func findAll(db aplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		users := make([]user, 0, len(db.data))

		for _, user := range db.data {
			users = append(users, user)
		}

		sendJSON(w, response{Data: users}, http.StatusOK)
	}
}

func handlePostUser(db *aplication) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var body user
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Error parsing the body of requisition", http.StatusUnprocessableEntity)
			return
		}

		uuid, err := uuid.NewUUID()

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		db.data[id(uuid)] = body

		p_response := postResponse{
			UUID: uuid.String(),
		}

		sendJSON(w, response{Data: p_response}, http.StatusCreated)

	}

}

func handlegetUser(db aplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		paramID := chi.URLParam(r, "id")

		ID, err := uuid.Parse(paramID)

		if err != nil {
			sendJSON(w, response{Message: "Something went wrong - Error parsing userID"}, http.StatusInternalServerError)
			return
		}

		user, ok := db.data[id(ID)]

		if !ok {
			sendJSON(w, response{Message: "User not found"}, http.StatusNotFound)
			return
		}

		sendJSON(w, response{Data: user}, http.StatusOK)

	}
}

func handleDeleteUser(db aplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		paramID := chi.URLParam(r, "id")

		ID, err := uuid.Parse(paramID)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		_, ok := db.data[id(ID)]

		if !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		delete(db.data, id(ID))

		sendJSON(w, response{Message: "User deleted"}, http.StatusOK)
	}
}

func putUser(db aplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		paramID := chi.URLParam(r, "id")

		ID, err := uuid.Parse(paramID)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		_, ok := db.data[id(ID)]

		if !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		var body user

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Error parsing the body of requisition", http.StatusUnprocessableEntity)
			return
		}

		db.data[id(ID)] = body

		sendJSON(w, response{Message: "User updated"}, http.StatusOK)
	}
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)

	})
}
