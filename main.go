package main

import (
	"encoding/json"
	"net/http"
)

type Todo struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul"`
	Deskripsi string `json:"deskripsi,omitempty"`
	Selesai   bool   `json:"selesai"`
}

var todo_list []Todo

func tambah(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Hanya mendukung method POST", http.StatusMethodNotAllowed)
		return
	}

	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Gagal decode JSON", http.StatusBadRequest)
		return
	}
	todo.ID = len(todo_list) + 1
	todo_list = append(todo_list, todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"pesan":  "Berhasil membuat data!",
		"data":   todo,
		"status": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

func lihatSemua(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Hanya mendukung method GET", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(todo_list)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/todo/tambah", tambah)
	mux.HandleFunc("/todo/lihat-semua", lihatSemua)

	http.ListenAndServe(":99", mux)
}
