package store

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// VENInfo represents the information of a VEN (Virtual Execution Node)
type VENInfo struct {
	Name                    string    `json:"name"`
	RAM                     string    `json:"ram"`
	Core                    string    `json:"core"`
	URL                     string    `json:"url"`
	MaxQueueSize            string    `json:"maxQueueSize"`
	CurrentQueueSize        string    `json:"currentQueueSize"`
	PreferenceList          string    `json:"preferenceList"`
	TrustScore              string    `json:"trustScore"`
	MaxQueueSizeUpdated     time.Time `json:"maxQueueSizeUpdated"`
	CurrentQueueSizeUpdated time.Time `json:"currentQueueSizeUpdated"`
	TrustScoreUpdated       time.Time `json:"trustScoreUpdated"`
}

func (s *MySQLStore) GetAvailableVEN() ([]VENInfo, error) {
	rows, err := s.db.Query("SELECT name, ram, core, url, max_queue_size, current_queue_size, preference_list, trust_score, max_queue_size_last_updated, current_queue_size_last_updated, trust_score_last_updated FROM ven_info where current_queue_size>0;")
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

func (s *MySQLStore) updateColumnAndTimestamp(venName string, columnName string, newValue string) error {
	var columnTimestampName string

	switch columnName {
	case "max_queue_size":
		columnTimestampName = "max_queue_size_last_updated"
	case "current_queue_size":
		columnTimestampName = "current_queue_size_last_updated"
	case "trust_score":
		columnTimestampName = "trust_score_last_updated"
	default:
		return errors.New("unsupported column name")
	}

	// Retrieve the current value of the column
	var currentValue string
	getQuery := fmt.Sprintf("SELECT %s FROM ven_info WHERE name = ?", columnName)
	err := s.db.QueryRow(getQuery, venName).Scan(&currentValue)
	if err != nil {
		return fmt.Errorf("failed to retrieve current value for %s: %v", venName, err)
	}

	// If the value hasn't changed, return without updating
	if currentValue == newValue {
		return nil
	}

	query := fmt.Sprintf("UPDATE ven_info SET %s = ?, %s = ? WHERE name = ?", columnName, columnTimestampName)
	_, err = s.db.Exec(query, newValue, time.Now(), venName)
	if err != nil {
		return fmt.Errorf("failed to execute update for %s: %v", venName, err)
	}
	log.Printf("Successfully updated queue size in DB for %s", venName)
	return nil
}

// UpdateMaxQueueSize updates the max_queue_size and its associated timestamp for a specific VENInfo.
func (s *MySQLStore) UpdateMaxQueueSize(venName string, newValue string) error {
	return s.updateColumnAndTimestamp(venName, "max_queue_size", newValue)
}

// UpdateCurrentQueueSize updates the current_queue_size and its associated timestamp for a specific VENInfo.
func (s *MySQLStore) UpdateCurrentQueueSize(venName string, newValue string) error {
	return s.updateColumnAndTimestamp(venName, "current_queue_size", newValue)
}

// UpdateTrustScore updates the trust_score and its associated timestamp for a specific VENInfo.
func (s *MySQLStore) UpdateTrustScore(venName string, newValue string) error {
	return s.updateColumnAndTimestamp(venName, "trust_score", newValue)
}

func (store *MySQLStore) CountVENInfo() (int, error) {
	var count int
	err := store.db.QueryRow("SELECT COUNT(*) FROM ven_info").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (store *MySQLStore) InsertVENInfo(info VENInfo) error {
	stmt, err := store.db.Prepare(`INSERT INTO ven_info 
        (name, url, ram, core, max_queue_size, current_queue_size, preference_list, trust_score, max_queue_size_last_updated, current_queue_size_last_updated, trust_score_last_updated) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(info.Name, info.URL, info.RAM, info.Core, info.MaxQueueSize, info.CurrentQueueSize, info.PreferenceList, info.TrustScore, info.MaxQueueSizeUpdated, info.CurrentQueueSizeUpdated, info.TrustScoreUpdated)
	return err
}
