package main

import (
	"encoding/json"
	"log"
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

var (
	todo_list []Todo
	lastId    int
)

// middleware logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log informasi request
		log.Printf("Method: %s, URL: %s, RemoteAddr: %s", r.Method, r.URL.Path, r.RemoteAddr)
		// lanjut ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

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
	lastId++
	todo.ID = lastId
	todo_list = append(todo_list, todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"pesan":  "Berhasil membuat data!",
		"data":   todo,
		"status": http.StatusCreated,
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

	response := map[string]interface{}{
		"pesan":  "Berhasil mengambil data!",
		"data":   todo_list,
		"status": http.StatusOK,
	}

	json.NewEncoder(w).Encode(response)
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
			return
		}
	}
	http.Error(w, "Todo tidak ditemukan!", http.StatusNotFound)
}

func perbarui(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Hanya mendukung method PUT", http.StatusMethodNotAllowed)
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
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "Todo tidak ditemukan", http.StatusNotFound)
}

func hapus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Hanya mendukung method DELETE", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todo/hapus/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid!", http.StatusBadRequest)
		return
	}

	for i, t := range todo_list {
		if t.ID == id {
			// Hapus dari slice
			todo_list = append(todo_list[:i], todo_list[i+1:]...)

			response := map[string]interface{}{
				"pesan":  "Berhasil menghapus data!",
				"status": http.StatusOK,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "Todo tidak ditemukan!", http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/todo/tambah", tambah)
	mux.HandleFunc("/todo/lihat-semua", lihatSemua)
	mux.HandleFunc("/todo/lihat-detail/", lihatDetail)
	mux.HandleFunc("/todo/perbarui/", perbarui)
	mux.HandleFunc("/todo/hapus/", hapus)

	// middleware logging
	loggedMux := loggingMiddleware(mux)

	http.ListenAndServe(":99", loggedMux)
}
