package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"scheduler-node/internal/store"
)

type Handler struct {
	workflowStore store.WorkflowStore
	venStore      store.VENStore
}

func NewHandler(wf store.WorkflowStore, store store.VENStore) *Handler {
	return &Handler{workflowStore: wf, venStore: store}
}

func (h *Handler) ShowTables(w http.ResponseWriter, r *http.Request) {
	// Get VEN info from the store
	// TODO remove pseudo methods
	venInfo, err := h.venStore.GetVENInfos()
	const layout = "2006-01-02 15:04:05"
	if err != nil {
		http.Error(w, "Failed to fetch VEN info", http.StatusInternalServerError)
		return
	}

	// Get Workflow info from the store
	workflowInfo, err := h.workflowStore.GetWorkflows(r.Context())
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
					<th>Submitted by</th>
					<th>Created at</th>
					<th>ram</th>
					<th>core</th>
					<th>policy</th>
					<th>Expected Execution time</th>
					<th>Actual Execution Time</th>
					<th>Assigned VM</th>
					<th>Processing_started_at</th>
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
				<td>` + workflow.SubmittedBy.String + `</td>
				<td>` + workflow.CreatedAt.Time.Local().Format(layout) + `</td>
				<td>` + workflow.RAM + `</td>
				<td>` + workflow.Core + `</td>
				<td>` + workflow.Policy + `</td>
				<td>` + workflow.ExpectedExecutionTime.String + `</td>
				<td>` + workflow.ActualExecutionTime.String + `</td>
				<td>` + workflow.AssignedVM.String + `</td>
				<td>` + workflow.ProcessingStartedAt.Time.Local().Format(layout) + `</td>
				<td>` + workflow.AssignedAt.Time.Local().Format(layout) + `</td>
				<td>` + workflow.CompletedAt.Time.Local().Format(layout) + `</td>
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

func stringToNullTime(s string) (sql.NullTime, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
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

// ...

type UpdateWorkflowInfo struct {
	Name               string `json:"Name"`
	Duration           string `json:"Duration"` // You might want to map this to a field in WorkflowInfo
	StartedExecutionAt string `json:"StartedExecutionAt"`
	CompletedAt        string `json:"CompletedAt"`
}

func stringToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func (h *Handler) UpdateWorkflowByName(w http.ResponseWriter, r *http.Request) {

	log.Println("UpdateWorkflowByName: Received request")
	// Decode the received JSON into the UpdateWorkflowInfo struct
	var updateInfo UpdateWorkflowInfo
	err := json.NewDecoder(r.Body).Decode(&updateInfo)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	log.Printf("UpdateWorkflowByName: Decoded JSON - %+v", updateInfo)
	startedExecutionAt, err := stringToNullTime(updateInfo.StartedExecutionAt)
	if err != nil {
		log.Printf("UpdateWorkflowByName: Error parsing StartedExecutionAt - %v", err)
		http.Error(w, "Invalid StartedExecutionAt format", http.StatusBadRequest)
		return
	}

	completedAt, err := stringToNullTime(updateInfo.CompletedAt)
	if err != nil {
		log.Printf("UpdateWorkflowByName: Error parsing CompletedAt - %v", err)
		http.Error(w, "Invalid CompletedAt format", http.StatusBadRequest)
		return
	}
	// Map the UpdateWorkflowInfo to WorkflowInfo and use the workflowStore
	// to update the workflow (assuming an UpdateWorkflow method exists)
	workflow := store.WorkflowInfo{
		Name:                updateInfo.Name,
		ActualExecutionTime: stringToNullString(updateInfo.Duration),
		// Map Duration to the appropriate field in WorkflowInfo if needed
		ProcessingStartedAt: startedExecutionAt,
		CompletedAt:         completedAt,
		Status:              store.WorkflowStatusDone,
	}

	err = h.workflowStore.UpdateWorkflowByName(r.Context(), &workflow)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Failed to update workflow", http.StatusInternalServerError)
		return
	}

	// Send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Workflow updated successfully"))
}
