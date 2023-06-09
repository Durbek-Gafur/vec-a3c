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
