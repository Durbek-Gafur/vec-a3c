package seeder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"scheduler-node/internal/store"
)

type Seeder interface {
	PopulateVENInfo() error
	PopulateWorkflows() error
	UpdateQueueSizePeriodically(ctx context.Context)
}
type dbSeeder struct {
	venStore    store.VENStore
	wfStore     store.WorkflowStore
	urlProvider URLProvider
}

func NewDBSeeder(venStore store.VENStore, wfStore store.WorkflowStore, urlProvider URLProvider) *dbSeeder {
	return &dbSeeder{venStore: venStore, wfStore: wfStore, urlProvider: urlProvider}
}

type ResourceSpec struct {
	RAM       string `json:"RAM"`
	CORE      string `json:"CORE"`
	MAX_QUEUE string `json:"MAX_QUEUE"`
}

type CurrentQueueSize struct {
	SIZE int `json:"size"`
}

func (ds *dbSeeder) PopulateVENInfo() error {
	// Check if the table is empty
	count, err := ds.venStore.CountVENInfo()
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

	for i := 1; i <= venCount; i++ {
		venName := "ven" + strconv.Itoa(i)
		url := ds.urlProvider.GetURL(venName)

		var venInfo store.VENInfo
		venInfo.Name = venName
		venInfo.URL = url
		venInfo.CurrentQueueSizeUpdated = time.Now()
		venInfo.MaxQueueSizeUpdated = time.Now()
		venInfo.TrustScoreUpdated = time.Now()

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

		if err := ds.venStore.InsertVENInfo(venInfo); err != nil {
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
	venInfo.TrustScore = fmt.Sprintf("%.2f", rand.Float64())

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
func (ds *dbSeeder) PopulateWorkflows() error {
	// Check if the table is empty
	count, err := ds.wfStore.CountWorkflows()
	if err != nil {
		return err
	}

	// If table is not empty, return
	if count > 0 {
		log.Println("Workflow info table is not empty, skipping population.")
		return nil
	}
	rand.Seed(time.Now().UnixNano())

	types := []string{"demo.fastq", "demo_25per.fastq", "demo_50per.fastq", "demo_75per.fastq"}
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
			wf := &store.WorkflowInfo{
				Name:        "workflow" + strconv.Itoa(i),
				Type:        types[rand.Intn(len(types))],
				RAM:         ramList[rand.Intn(len(ramList))],
				Core:        coreList[rand.Intn(len(coreList))],
				Policy:      policies[rand.Intn(len(policies))],
				SubmittedBy: sql.NullString{String: userList[rand.Intn(len(userList))], Valid: true},
				CreatedAt:   arrivalTimes[i-1],
			}
			err := ds.wfStore.InsertWorkflow(wf)

			if err != nil {
				log.Printf("Failed to populate workflow: %v", err)
				return err
			}
		}
	}

	return nil
}

func generateArrivalTimes(startTime, endTime, lambda float64) []sql.NullTime {
	rand.Seed(2023)
	interval := 1.0 // in minutes

	// Convert start and end time to minutes
	startMinutes := startTime * 60
	endMinutes := endTime * 60

	var arrivalTimes []sql.NullTime

	// Generate sequence of inter-arrival times based on exponential distribution
	for currentTime := startMinutes; currentTime <= endMinutes; {
		// Generate exponentially distributed random number
		interArrivalTime := rand.ExpFloat64() / lambda
		currentTime += interArrivalTime * interval

		if currentTime <= endMinutes {
			// Convert to sql.NullTime and append to arrivalTimes
			arrivalTime := sql.NullTime{
				Time:  time.Unix(int64(currentTime*60), 0), // convert minutes to seconds
				Valid: true,
			}
			arrivalTimes = append(arrivalTimes, arrivalTime)
		}
	}

	return arrivalTimes
}

func (ds *dbSeeder) UpdateQueueSizePeriodically(ctx context.Context) {
	log.Println("Starting UpdateQueueSizePeriodically...")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Generate or fetch VENs
	venInfos := ds.fetchOrDefineVENs()
	log.Printf("Fetched or defined %d VENs", len(venInfos))

	for {
		select {
		case <-ticker.C:
			start := time.Now()
			log.Println("New tick started")

			// Using a WaitGroup to ensure all parallel tasks are completed
			var wg sync.WaitGroup
			for _, venInfo := range venInfos {
				wg.Add(1) // Increment the counter before starting the goroutine
				log.Printf("Initiating fetch for %s", venInfo.URL)

				go func(venInfo store.VENInfo) { // Note: Passing `venInfo` as a parameter to avoid data race
					defer wg.Done() // Decrement the counter once the goroutine finishes

					var retryCount int
					var queueSize string
					var err error

					// Retry fetching the queue size up to 3 times with 1-second intervals.
					for retryCount < 3 {
						queueSize, err = FetchQueueSize(venInfo.URL + "/queue-size")
						if err == nil {
							break
						}
						log.Printf("Retry %d failed for %s: %v", retryCount+1, venInfo.URL, err)
						retryCount++
						time.Sleep(1 * time.Second)
					}

					if err != nil {
						log.Printf("Failed to fetch queue size for %s after retries: %v", venInfo.URL, err)
						return
					}

					err = ds.venStore.UpdateCurrentQueueSize(venInfo.Name, queueSize)
					if err != nil {
						log.Printf("Failed to update queue size in DB for %s: %v", venInfo.Name, err)
					}
				}(venInfo) // Passing current `venInfo` as argument
			}

			wg.Wait() // Wait until all goroutines are finished
			log.Println("All goroutines completed for this tick")

			// Check elapsed time to avoid overlapping ticks.
			elapsed := time.Since(start)
			if elapsed >= 5*time.Second {
				log.Println("Skipped tick due to overlap")
			} else {
				log.Printf("Tick completed in %v", elapsed)
			}

		case <-ctx.Done():
			log.Println("Context is done. Exiting UpdateQueueSizePeriodically...")
			return
		}
	}
}

// fetchOrDefineVENs generates a list of VENInfo based on the VEN_COUNT environment variable.
func (ds *dbSeeder) fetchOrDefineVENs() []store.VENInfo {
	venCount, err := strconv.Atoi(os.Getenv("VEN_COUNT"))
	if err != nil {
		log.Fatalf("Failed to fetch VEN_COUNT: %v", err)
	}

	var venInfos []store.VENInfo

	for i := 1; i <= venCount; i++ {
		venName := "ven" + strconv.Itoa(i)
		url := ds.urlProvider.GetURL(venName)

		var venInfo store.VENInfo
		venInfo.Name = venName
		venInfo.URL = url

		venInfos = append(venInfos, venInfo)
	}

	return venInfos
}
