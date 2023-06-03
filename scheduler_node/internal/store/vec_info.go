package store

import (
	"time"
)

// VENInfo represents the information of a VEN (Virtual Execution Node)
type VENInfo struct {
	Name              string `json:"name"`
	RAM               string `json:"ram"`
	Core              string `json:"core"`
	URL               string `json:"url"`
	MaxQueueSize      string `json:"maxQueueSize"`
	CurrentQueueSize  string `json:"currentQueueSize"`
	PreferenceList    string `json:"preferenceList"`
	TrustScore        string `json:"trustScore"`
	MaxQueueSizeUpdated   time.Time `json:"maxQueueSizeUpdated"`
	CurrentQueueSizeUpdated time.Time `json:"currentQueueSizeUpdated"`
	TrustScoreUpdated     time.Time `json:"trustScoreUpdated"`
}


// GetVENInfo retrieves VEN information from the database
func (s *MySQLStore) GetVENInfo() ([]VENInfo, error) {
	// TODO: Implement the actual retrieval of VEN info from the database
	// For this example, let's assume the data is retrieved from the database successfully

	// Simulated pseudo data in case of error
	pseudoData := []VENInfo{
		{
			Name:            "Ven 1",
			RAM:             "8GB",
			Core:            "4",
			URL:             "https://ven1.example.com",
			MaxQueueSize:    "100",
			CurrentQueueSize: "10",
			PreferenceList:  "User A",
			TrustScore:      "0.9",
			MaxQueueSizeUpdated:  time.Now(),
            CurrentQueueSizeUpdated: time.Now().Add(-time.Hour*25),
            TrustScoreUpdated:    time.Now(),
		},
		{
			Name:            "Ven 2",
			RAM:             "16GB",
			Core:            "8",
			URL:             "https://ven2.example.com",
			MaxQueueSize:    "200",
			CurrentQueueSize: "20",
			PreferenceList:  "User A, User B",
			TrustScore:      "0.8",
			MaxQueueSizeUpdated:  time.Now(),
            CurrentQueueSizeUpdated: time.Now(),
            TrustScoreUpdated:    time.Now().Add(-time.Minute*8),
		},
	}

	// Return the pseudo data in case of error
	return pseudoData, nil
}


func (s *MySQLStore) GetVENInfos() ([]VENInfo, error) {
    rows, err := s.db.Query("SELECT name, ram, core, url, max_queue_size, current_queue_size, preference_list, trust_score, max_queue_size_last_updated, current_queue_size_last_updated, trust_score_last_updated FROM ven_info")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var venInfos []VENInfo
    for rows.Next() {
        var venInfo VENInfo
        err := rows.Scan(&venInfo.Name, &venInfo.RAM, &venInfo.Core, &venInfo.URL, &venInfo.MaxQueueSize, &venInfo.CurrentQueueSize, &venInfo.PreferenceList, &venInfo.TrustScore, &venInfo.MaxQueueSizeUpdated, &venInfo.CurrentQueueSizeUpdated, &venInfo.TrustScoreUpdated)
        if err != nil {
            return nil, err
        }
        venInfos = append(venInfos, venInfo)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return venInfos, nil
}
