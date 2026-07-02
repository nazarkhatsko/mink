package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const apiKey = "secret-api-key"

type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	ID     int     `json:"id"`
	Status string  `json:"status"`
	Items  []Item  `json:"items"`
	Total  float64 `json:"total"`
}

var (
	mu      sync.Mutex
	orders  = map[int]*Order{}
	nextOID = 1
	nextIID = 1
)

func auth(r *http.Request) bool {
	return r.Header.Get("X-Api-Key") == apiKey
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	mux := http.NewServeMux()

	// POST /orders
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if !auth(r) {
			writeJSON(w, 401, map[string]string{"error": "unauthorized"})
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			return
		}
		mu.Lock()
		o := &Order{ID: nextOID, Status: "pending", Items: []Item{}}
		orders[nextOID] = o
		nextOID++
		mu.Unlock()
		writeJSON(w, 201, map[string]any{"data": o})
	})

	// POST /orders/:id/items  and  GET /orders/:id  and  POST /orders/:id/checkout
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if !auth(r) {
			writeJSON(w, 401, map[string]string{"error": "unauthorized"})
			return
		}

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		// parts: ["orders", "<id>"] or ["orders", "<id>", "items"|"checkout"]

		if len(parts) < 2 {
			w.WriteHeader(400)
			return
		}

		oid, err := strconv.Atoi(parts[1])
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid order id"})
			return
		}

		mu.Lock()
		o, ok := orders[oid]
		mu.Unlock()
		if !ok {
			writeJSON(w, 404, map[string]string{"error": "order not found"})
			return
		}

		// GET /orders/:id
		if len(parts) == 2 && r.Method == http.MethodGet {
			mu.Lock()
			writeJSON(w, 200, map[string]any{"data": o})
			mu.Unlock()
			return
		}

		// POST /orders/:id/items
		if len(parts) == 3 && parts[2] == "items" && r.Method == http.MethodPost {
			var body struct {
				Name  string  `json:"name"`
				Price float64 `json:"price"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
				writeJSON(w, 400, map[string]string{"error": "invalid body"})
				return
			}
			if body.Price == 0 {
				body.Price = float64(rand.Intn(9000)+1000) / 100.0
			}
			mu.Lock()
			item := Item{ID: nextIID, Name: body.Name, Price: body.Price}
			nextIID++
			o.Items = append(o.Items, item)
			o.Total = 0
			for _, it := range o.Items {
				o.Total += it.Price
			}
			mu.Unlock()
			writeJSON(w, 201, map[string]any{"data": item})
			return
		}

		// POST /orders/:id/checkout
		if len(parts) == 3 && parts[2] == "checkout" && r.Method == http.MethodPost {
			mu.Lock()
			if o.Status == "completed" {
				mu.Unlock()
				writeJSON(w, 409, map[string]string{"error": "already checked out"})
				return
			}
			o.Status = "completed"
			snap := *o
			mu.Unlock()
			writeJSON(w, 200, map[string]any{"data": snap})
			return
		}

		w.WriteHeader(404)
	})

	fmt.Println("order-flow server listening on :8081")
	http.ListenAndServe(":8081", mux)
}
