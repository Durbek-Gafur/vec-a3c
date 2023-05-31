package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"scheduler-node/internal/store"

	"github.com/gorilla/mux"
)

type Handler struct {
	queueSizeStore store.QueueSizeStore
	workflowStore store.WorkflowStore
	store store.Store
}

func NewHandler(qs store.QueueSizeStore,wf store.WorkflowStore, store store.Store) *Handler {
	return &Handler{queueSizeStore: qs, workflowStore: wf, store: store}
}

func (h *Handler) GetQueueSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	size, err := h.queueSizeStore.GetQueueSize(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error. Failed to fetch queue size", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]int{"size": size}, http.StatusOK)
	
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

	err = h.queueSizeStore.SetQueueSize(ctx, intSize)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]int{"size": intSize}, http.StatusOK)
	
	
}

func (h *Handler) UpdateQueueSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
		fmt.Println("Invalid Request couldn't decode")
		http.Error(w, "Invalid Request couldn't decode", http.StatusBadRequest)
		return
	}

	size, ok := data["size"]
	
	if !ok {
		fmt.Println("Invalid Request no 'size'")
		http.Error(w, "Invalid Request no 'size'", http.StatusBadRequest)
		return
	}

	intSize, err := strconv.Atoi(size)
	if err != nil {
		fmt.Println("Invalid Request 'size' should be an integer")
		http.Error(w, "Invalid Request 'size' should be an integer", http.StatusBadRequest)
		return
	}

	err = h.queueSizeStore.UpdateQueueSize(ctx, intSize)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, map[string]string{"size": size}, http.StatusOK)
	
	
}

func (h *Handler) GetWorkflowByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	wf, err := h.workflowStore.GetWorkflowByID(ctx,int(id))
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

	workflows, err := h.workflowStore.GetWorkflows(ctx,filter)
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

	if _,err := h.workflowStore.SaveWorkflow(ctx,&wf); err != nil {
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

func (h *Handler) ShowTables(w http.ResponseWriter, r *http.Request) {
	// Get VEN info from the store
	venInfo, err := h.store.GetVENInfo()
	if err != nil {
		http.Error(w, "Failed to fetch VEN info", http.StatusInternalServerError)
		return
	}

	// Get Workflow info from the store
	workflowInfo, err := h.store.GetWorkflowInfo()
	if err != nil {
		http.Error(w, "Failed to fetch Workflow info", http.StatusInternalServerError)
		return
	}

	// Generate HTML response
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Scheduler Server's Tables</title>
			<style>
			table {
				border-collapse: collapse;
				width: 100%;
			}
			
			th, td {
				border: 1px solid black;
				padding: 8px;
			}
		</style>
		</head>
		<body>
			<h1>VEN info</h1>
			<table>
				<tr>
					<th>VEN name</th>
					<th>ram</th>
					<th>core</th>
					<th>url</th>
					<th>Max queue size</th>
					<th>Current queue size</th>
					<th>Preference list</th>
					<th>Trust Score</th>
				</tr>
	`

	// Generate rows for VEN info table
	for _, ven := range venInfo {
		maxQueueSizeUpdated := time.Since(ven.MaxQueueSizeUpdated).Round(time.Second)
		currentQueueSizeUpdated := time.Since(ven.CurrentQueueSizeUpdated).Round(time.Second)
		trustScoreUpdated := time.Since(ven.TrustScoreUpdated).Round(time.Second)

		html += `
			<tr>
				<td>` + ven.Name + `</td>
				<td>` + ven.RAM + `</td>
				<td>` + ven.Core + `</td>
				<td>` + ven.URL + `</td>
				<td>` + ven.MaxQueueSize + `<br>Updated ` + formatTimeAgo(maxQueueSizeUpdated) + `</td>
				<td>` + ven.CurrentQueueSize + `<br>Updated ` + formatTimeAgo(currentQueueSizeUpdated) + `</td>
				<td>` + ven.PreferenceList + `</td>
				<td>` + ven.TrustScore + `<br>Updated ` + formatTimeAgo(trustScoreUpdated) + `</td>
			</tr>
		`
	}


	html += `
			</table>
			<h1>Workflow info</h1>
			<table>
				<tr>
					<th>Workflow Name</th>
					<th>Type</th>
					<th>ram</th>
					<th>core</th>
					<th>policy</th>
					<th>Expected Execution time</th>
					<th>Actual Execution Time</th>
					<th>Assigned VM</th>
					<th>Assigned_at</th>
					<th>Completed_at</th>
					<th>Status</th>
				</tr>
	`

	// Generate rows for Workflow info table
	for _, workflow := range workflowInfo {
		lastUpdated := time.Since(workflow.LastUpdated).Round(time.Second)
		html += `
			<tr>
				<td>` + workflow.Name + `</td>
				<td>` + workflow.Type + `</td>
				<td>` + workflow.RAM + `</td>
				<td>` + workflow.Core + `</td>
				<td>` + workflow.Policy + `</td>
				<td>` + workflow.ExpectedExecutionTime + `</td>
				<td>` + workflow.ActualExecutionTime + `</td>
				<td>` + workflow.AssignedVM + `</td>
				<td>` + workflow.AssignedAt + `</td>
				<td>` + workflow.CompletedAt + `</td>
				<td>` + workflow.Status + `<br>Updated ` + formatTimeAgo(lastUpdated) + `</td>
			</tr>
		`
	}

	html += `
			</table>
		</body>
		</html>
	`

	// Set the Content-Type header to indicate a HTML response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Write the HTML response to the client
	w.Write([]byte(html))
}


// formatTimeAgo formats the time duration as "n times ago"
func formatTimeAgo(duration time.Duration) string {
    seconds := int(duration.Seconds())

    if seconds == 0 {
        return "just now"
    } else if seconds < 60 {
        return fmt.Sprintf("%d seconds ago", seconds)
    } else if seconds < 3600 {
        minutes := seconds / 60
        return fmt.Sprintf("%d minutes ago", minutes)
    } else if seconds < 86400 {
        hours := seconds / 3600
        return fmt.Sprintf("%d hours ago", hours)
    } else {
        days := seconds / 86400
        return fmt.Sprintf("%d days ago", days)
    }
}
