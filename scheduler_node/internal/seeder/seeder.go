package seeder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"scheduler-node/internal/store"
)

type ResourceSpec struct {
	RAM        string `json:"RAM"`
	CORE       string `json:"CORE"`
	MAX_QUEUE  string `json:"MAX_QUEUE"`
}

type CurrentQueueSize struct {
	SIZE  int `json:"size"`
}

func PopulateVENInfo(db *sql.DB,urlProvider URLProvider) error {
	// Check if the table is empty
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM ven_info").Scan(&count)
	if err != nil {
		return err
	}

	// If table is not empty, return
	if count > 0 {
		log.Println("VEN info table is not empty, skipping population.")
		return nil
	}

	venCount, err := strconv.Atoi(os.Getenv("VEN_COUNT"))
	if err != nil {
		return fmt.Errorf("invalid VEN_COUNT: %w", err)
	}

	stmt, err := db.Prepare(`INSERT INTO ven_info 
    (name, url, ram, core, max_queue_size, current_queue_size, preference_list, trust_score, max_queue_size_last_updated, current_queue_size_last_updated, trust_score_last_updated) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 1; i <= venCount; i++ {
		venName := "ven" + strconv.Itoa(i)
		url := urlProvider.GetURL(venName)

		var venInfo store.VENInfo
		venInfo.Name = venName
		venInfo.URL = url

		venInfoChan := make(chan store.VENInfo)
		errChan := make(chan error)

		go func() {
			var err error
			for retries := 0; retries < 5; retries++ {
				if err = FetchAndPopulateData(&venInfo); err != nil {
					time.Sleep(time.Second * time.Duration(retries*2))
					continue
				}
				venInfoChan <- venInfo
				return
			}
			errChan <- err
		}()

		select {
		case venInfo = <-venInfoChan:
		case err = <-errChan:
			return err
		}

		preferenceListStr, err := json.Marshal(venInfo.PreferenceList)
		if err != nil {
			return fmt.Errorf("failed to marshal preference list: %w", err)
		}

		if _, err := stmt.Exec(venInfo.Name, venInfo.URL, venInfo.RAM, venInfo.Core, venInfo.MaxQueueSize, venInfo.CurrentQueueSize, string(preferenceListStr), venInfo.TrustScore, time.Now(), time.Now(), time.Now()); err != nil {
			return err
		}
	}
	return nil
}


func FetchAndPopulateData(venInfo *store.VENInfo) error {
	// Fetching RAM, CORE and MAX_QUEUE from url/rspec
	resSpec, err := FetchResourceSpec(venInfo.URL + "/rspec")
	if err != nil {
		return err
	}
	venInfo.RAM = resSpec.RAM
	venInfo.Core = resSpec.CORE
	venInfo.MaxQueueSize = resSpec.MAX_QUEUE

	// Fetching current queue size from url/queue-size
	queueSize, err := FetchQueueSize(venInfo.URL + "/queue-size")
	if err != nil {
		return err
	}
	venInfo.CurrentQueueSize = queueSize

	// Generating random Preference List
	venInfo.PreferenceList = GeneratePreferenceList()

	// Generating random Trust Score
	venInfo.TrustScore = strconv.FormatFloat(rand.Float64(), 'e', -1, 64)

	return nil
}

func FetchResourceSpec(url string) (*ResourceSpec, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status: got %v", res.StatusCode)
	}

	var resSpec ResourceSpec
	if err := json.NewDecoder(res.Body).Decode(&resSpec); err != nil {
		return nil, err
	}
	return &resSpec, nil
}

func FetchQueueSize(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected http status: got %v", res.StatusCode)
	}
	var current_queue_size CurrentQueueSize
	if err := json.NewDecoder(res.Body).Decode(&current_queue_size); err != nil {
		return "", err
	}
	return strconv.Itoa(current_queue_size.SIZE), nil
}

func GeneratePreferenceList() string {
	userList := []string{"UserA", "UserB", "UserC", "UserD", "UserE", "UserF", "UserG", "UserH", "UserI", "UserJ"}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(userList), func(i, j int) { userList[i], userList[j] = userList[j], userList[i] })
	numElements := rand.Intn(len(userList)-1) + 1
	selected := userList[:numElements]
	return strings.Join(selected, ", ")
}


// PopulateWorkflows populates the workflow_info table with the provided workflows.
func PopulateWorkflows(db *sql.DB) error {
	// Check if the table is empty
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM workflow_info").Scan(&count)
	if err != nil {
		return err
	}

	// If table is not empty, return
	if count > 0 {
		log.Println("Workflow info table is not empty, skipping population.")
		return nil
	}
	rand.Seed(time.Now().UnixNano())

	types := []string{"typeA", "typeB"}
	ramList := []string{"512Mi", "1Gi", "1.5Gi", "2Gi", "2.5Gi", "3Gi"} 
	coreList := []string{"0.5", "1", "1.5", "2", "2.5", "3"} 
	userList := []string{"UserA", "UserB", "UserC", "UserD", "UserE", "UserF", "UserG", "UserH", "UserI", "UserJ"}
	policies := []string{"policyA", "policyB", "policyC"} 

	maxWf, err := strconv.Atoi(os.Getenv("MAX_WF"))
	if err != nil {
		log.Printf("Failed to parse MAX_WF: %v", err)
		return err
	}

	// Generate arrival times
	arrivalTimes := generateArrivalTimes(8, 17, 0.5) // start time, end time and lambda

	for i := 1; i <= maxWf; i++ {
		if i-1 < len(arrivalTimes) {
			createdAt := arrivalTimes[i-1]

			_, err := db.Exec(`
				INSERT INTO workflow_info 
				(name, type, ram, core, policy, submitted_by, created_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				"workflow" + strconv.Itoa(i), // generate name
				types[rand.Intn(len(types))], // pick randomly from types
				ramList[rand.Intn(len(ramList))], // pick randomly from ramList
				coreList[rand.Intn(len(coreList))], // pick randomly from coreList
				policies[rand.Intn(len(policies))], // pick randomly from policies
				userList[rand.Intn(len(userList))], // pick randomly from userList
				createdAt, // use generated createdAt
			)

			if err != nil {
				log.Printf("Failed to populate workflow: %v", err)
				return err
			}
		}
	}

	return nil
}


func generateArrivalTimes(startTime, endTime, lambda float64) []time.Time {
	rand.Seed(2023)
	interval := 1.0 // in minutes

	// Convert start and end time to minutes
	startMinutes := startTime * 60
	endMinutes := endTime * 60

	var arrivalTimes []time.Time

	// Generate sequence of inter-arrival times based on exponential distribution
	for currentTime := startMinutes; currentTime <= endMinutes; {
		// Generate exponentially distributed random number
		interArrivalTime := rand.ExpFloat64() / lambda
		currentTime += interArrivalTime * interval

		if currentTime <= endMinutes {
			// Convert to time.Time and append to arrivalTimes
			arrivalTime := time.Unix(int64(currentTime*60), 0) // convert minutes to seconds
			arrivalTimes = append(arrivalTimes, arrivalTime)
		}
	}

	return arrivalTimes
}
