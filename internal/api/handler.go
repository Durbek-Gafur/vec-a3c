package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"vec-node/internal/store"

	"github.com/gorilla/mux"
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
		http.Error(w, "Internal Server Error. Failed to fetch queue size", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]int{"size": size}, http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"size":` + strconv.Itoa(size) + `}`))
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

	jsonResponse(w, map[string]int{"size": intSize}, http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"size":` + size + `}`))
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

	jsonResponse(w, map[string]string{"size": size}, http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"size":` + size + `}`))
}

func (h *Handler) GetWorkflowByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	wf, err := h.store.GetWorkflowByID(ctx,id)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, wf, http.StatusOK)
}

func (h *Handler) GetWorkflows(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filter := &store.WorkflowFilter{
		Type:      r.URL.Query().Get("type"),
		StartTime: parseTime(r.URL.Query().Get("start_time")),
		EndTime:   parseTime(r.URL.Query().Get("end_time")),
	}

	workflows, err := h.store.GetWorkflows(ctx,filter)
	if err != nil {
		http.Error(w, "Failed to fetch workflows", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, workflows, http.StatusOK)
}

func (h *Handler) SaveWorkflow(w http.ResponseWriter, r *http.Request) {
	var wf store.Workflow
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _,err := h.store.SaveWorkflow(ctx,&wf); err != nil {
		http.Error(w, "Failed to save workflow", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, wf, http.StatusCreated)
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}

	return t
}