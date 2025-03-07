package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

func lihatDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Hanya mendukung method GET", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todo/lihat-detail/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid!", http.StatusBadRequest)
		return
	}

	for _, i := range todo_list {
		if id == i.ID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(i)
		}
	}
}

func perbarui(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Hanya mendukung method Put", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todo/perbarui/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid!", http.StatusBadRequest)
		return
	}

	var updateTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updateTodo)
	if err != nil {
		http.Error(w, "Gagal decode JSON", http.StatusBadRequest)
		return
	}

	for i, t := range todo_list {
		if t.ID == id {
			todo_list[i].Judul = updateTodo.Judul
			todo_list[i].Deskripsi = updateTodo.Deskripsi
			todo_list[i].Selesai = updateTodo.Selesai

			response := map[string]interface{}{
				"pesan":  "Berhasil memperbarui data!",
				"data":   todo_list[i],
				"status": http.StatusOK,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "Todo tidak ditemukan", http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/todo/tambah", tambah)
	mux.HandleFunc("/todo/lihat-semua", lihatSemua)
	mux.HandleFunc("/todo/lihat-detail/{id}", lihatDetail)
	mux.HandleFunc("/todo/perbarui/{id}", perbarui)

	http.ListenAndServe(":99", mux)
}
