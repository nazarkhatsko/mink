package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const apiKey = "secret-api-key"

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var (
	mu    sync.RWMutex
	users = map[string]User{}
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", auth(createUser))
	mux.HandleFunc("GET /users", auth(listUsers))
	mux.HandleFunc("GET /users/{id}", auth(getUser))
	mux.HandleFunc("DELETE /users/{id}", auth(deleteUser))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Api-Key")
		if key != apiKey {
			writeError(w, http.StatusUnauthorized, "invalid api key")
			return
		}
		next(w, r)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if u.Name == "" || u.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email are required")
		return
	}

	u.ID = generateID()

	mu.Lock()
	users[u.ID] = u
	mu.Unlock()

	writeJSON(w, http.StatusCreated, map[string]any{"data": u})
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	list := make([]User, 0, len(users))
	for _, u := range users {
		list = append(list, u)
	}
	mu.RUnlock()

	writeJSON(w, http.StatusOK, map[string]any{"data": list, "count": len(list)})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	mu.RLock()
	u, ok := users[id]
	mu.RUnlock()

	if !ok {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": u})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	mu.Lock()
	_, ok := users[id]
	delete(users, id)
	mu.Unlock()

	if !ok {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"message": "deleted"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}

var (
	idCounter int
	idMu      sync.Mutex
)

func generateID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return fmt.Sprintf("%04d", idCounter)
}
