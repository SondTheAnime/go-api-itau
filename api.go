package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserService interface {
	GetUsers() ([]User, error)
	CreateUser(user User) (User, error)
}

type InMemoryUserService struct {
	users []User
}

func (s *InMemoryUserService) GetUsers() ([]User, error) {
	return s.users, nil
}

func (s *InMemoryUserService) CreateUser(user User) (User, error) {
	s.users = append(s.users, user)
	return user, nil
}

type UserHandler struct {
	service UserService
}

func (h *UserHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetUsers(w, r)
	case http.MethodPost:
		h.handleCreateUser(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdUser, err := h.service.CreateUser(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	userService := &InMemoryUserService{
		users: []User{
			{ID: "1", Name: "João"},
			{ID: "2", Name: "Maria"},
		},
	}
	userHandler := &UserHandler{
		service: userService,
	}

	router.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.handleGetUsers)
		r.Post("/", userHandler.handleCreateUser)
	})

	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "API de Usuários - Documentação",
			},
			DarkMode: true,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, htmlContent)
	})

	log.Println("Servidor iniciado em http://localhost:8080")
	log.Println("Documentação disponível em http://localhost:8080/docs")
	log.Fatal(http.ListenAndServe(":8080", router))
}
