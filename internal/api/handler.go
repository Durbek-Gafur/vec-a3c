package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"vec-node/internal/store"
)

type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) GetQueueSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("hello from GET")
	size, err := h.store.GetQueueSize(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"size":` + strconv.Itoa(size) + `}`))
}

func (h *Handler) SetQueueSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request couldn't decode", http.StatusBadRequest)
		return
	}

	size, ok := data["size"]
	if !ok {
		http.Error(w, "Invalid Request no 'size'", http.StatusBadRequest)
		return
	}

	intSize, err := strconv.Atoi(size)
	if err != nil {
		http.Error(w, "Invalid Request 'size' should be an integer", http.StatusBadRequest)
		return
	}

	err = h.store.SetQueueSize(ctx, intSize)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"size":` + size + `}`))
}

func (h *Handler) UpdateQueueSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request couldn't decode", http.StatusBadRequest)
		return
	}

	size, ok := data["size"]
	if !ok {
		http.Error(w, "Invalid Request no 'size'", http.StatusBadRequest)
		return
	}

	intSize, err := strconv.Atoi(size)
	if err != nil {
		http.Error(w, "Invalid Request 'size' should be an integer", http.StatusBadRequest)
		return
	}

	err = h.store.UpdateQueueSize(ctx, intSize)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"size":` + size + `}`))
}
