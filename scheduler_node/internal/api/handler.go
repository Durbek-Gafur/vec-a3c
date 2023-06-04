package api

import (
	"fmt"
	"net/http"
	"time"

	"scheduler-node/internal/store"
)

type Handler struct {
	workflowStore store.WorkflowStore
	venStore store.VENStore
}

func NewHandler(wf store.WorkflowStore, store store.VENStore) *Handler {
	return &Handler{ workflowStore: wf, venStore: store}
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
	workflowInfo, err := h.workflowStore.PseudoGetWorkflowInfo(r.Context())
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
				<td>` + workflow.SubmittedBy + `</td>
				<td>` + workflow.RAM + `</td>
				<td>` + workflow.Core + `</td>
				<td>` + workflow.Policy + `</td>
				<td>` + workflow.ExpectedExecutionTime + `</td>
				<td>` + workflow.ActualExecutionTime + `</td>
				<td>` + workflow.AssignedVM + `</td>
				<td>` + workflow.AssignedAt.Local().Format(layout) + `</td>
				<td>` + workflow.CompletedAt.Local().Format(layout) + `</td>
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
