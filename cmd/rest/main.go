// REST API for Catalog bounded context (Skillab HW5 task 01).
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/internal/catalog"
)

func main() {
	store := catalog.NewStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, store.List())
	})
	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/products/")
		if p, ok := store.Get(id); ok {
			writeJSON(w, p)
			return
		}
		http.NotFound(w, r)
	})

	addr := ":8080"
	log.Printf("catalog REST listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
