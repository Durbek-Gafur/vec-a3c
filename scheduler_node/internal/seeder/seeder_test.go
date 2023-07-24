package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"scheduler-node/internal/store"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)


func TestGeneratePreferenceList(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	expectedUserList := []string{"UserA", "UserB", "UserC", "UserD", "UserE", "UserF", "UserG", "UserH", "UserI", "UserJ"}

	// Call the function multiple times and check if the generated list is valid
	for i := 0; i < 10; i++ {
		result := GeneratePreferenceList()
		resultList := strings.Split(result, ", ")

		// Check if the number of elements in the generated list is within the expected range
		if len(resultList) < 1 || len(resultList) > len(expectedUserList) {
			t.Errorf("Generated list has an invalid number of elements: %s", result)
		}

		// Check if all the elements in the generated list are from the expected user list
		for _, user := range resultList {
			if !contains(expectedUserList, user) {
				t.Errorf("Generated list contains an invalid user: %s", user)
			}
		}
	}
}

// Helper function to check if a slice contains a specific element
func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}




func TestFetchAndPopulateData(t *testing.T) {

	mockServer := setupMockServer()
	defer mockServer.Close()

	venInfo := &store.VENInfo{
		Name: "ven1",
		URL:  mockServer.URL,
	}

	err := FetchAndPopulateData(venInfo)

	

	assert.NoError(t, err)
	assert.Equal(t, "600mi", venInfo.RAM)
	assert.Equal(t, "0.5", venInfo.Core)
	assert.Equal(t, "7", venInfo.MaxQueueSize)
	assert.Equal(t, "3", venInfo.CurrentQueueSize)
	assert.NotEmpty(t, venInfo.PreferenceList)

	trustScore, _ := strconv.ParseFloat(venInfo.TrustScore, 64)
	if trustScore < 0 || trustScore > 1 {
		t.Errorf("Trust score is out of range, got: %s", venInfo.TrustScore)
	}
}

func setupMockServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/rspec", mockRspecHandler)
	handler.HandleFunc("/queue-size", mockQueueSizeHandler)

	srv := httptest.NewServer(handler)

	return srv
}

func mockRspecHandler(w http.ResponseWriter, r *http.Request) {
	resp := `{"RAM":"600mi", "CORE":"0.5", "MAX_QUEUE":"7"}`
	fmt.Fprintln(w, resp)
}

func mockQueueSizeHandler(w http.ResponseWriter, r *http.Request) {
	resp := `{"size":3}`
	fmt.Fprintln(w, resp)
}

func TestFetchVENResources(t *testing.T) {
	mockServer := setupMockServer()
	defer mockServer.Close()

	venrspec,err := FetchResourceSpec(mockServer.URL+"/rspec")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if venrspec.RAM != "600mi" || venrspec.CORE != "0.5" || venrspec.MAX_QUEUE != "7" {
		t.Errorf("Unexpected resources data")
	}
}

func TestFetchCurrentQueueSize(t *testing.T) {
	mockServer := setupMockServer()
	defer mockServer.Close()


	queueSize,err := FetchQueueSize(mockServer.URL+"/queue-size")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if queueSize != "3" {
		t.Errorf("Unexpected queue size data")
	}
}


func TestPopulateVENInfo(t *testing.T) {
	os.Setenv("VEN_COUNT", "2")

	// Mock server for HTTP calls
	mockServer := setupMockServer()
	defer mockServer.Close()

	mockUrlProvider := NewMockURLProvider(mockServer.URL)

	// Mock database using sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expect QueryRow statement for checking if table is empty
	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT(.*) FROM ven_info").WillReturnRows(rows)


	// Expect Prepare statement
	mock.ExpectPrepare("INSERT INTO ven_info")

	// Expect Exec calls
	for i := 1; i <= 2; i++ {
		mock.ExpectExec("INSERT INTO ven_info").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(int64(i), 1))
	}

	// Call function under test
	err = PopulateVENInfo(db, mockUrlProvider)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPopulateVENInfoWithNonEmptyDB(t *testing.T) {
	os.Setenv("VEN_COUNT", "2")

	// Mock server for HTTP calls
	mockServer := setupMockServer()
	defer mockServer.Close()

	mockUrlProvider := NewMockURLProvider(mockServer.URL)

	// Mock database using sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expect QueryRow statement for checking if table is empty
	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT(.*) FROM ven_info").WillReturnRows(rows)

	// Call function under test
	err = PopulateVENInfo(db, mockUrlProvider)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPopulateWorkflows(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set MAX_WF environment variable
	os.Setenv("MAX_WF", "10")

	// Set expectation for SELECT COUNT query
	mock.ExpectQuery("SELECT COUNT(.*) FROM workflow_info").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Set expectations for exec
	for i := 1; i <= 10; i++ {
		mock.ExpectExec("INSERT INTO workflow_info").
			WithArgs(
				sqlmock.AnyArg(), // name
				sqlmock.AnyArg(), // type
				sqlmock.AnyArg(), // ram
				sqlmock.AnyArg(), // core
				sqlmock.AnyArg(), // policy
				sqlmock.AnyArg(), // submitted_by
				sqlmock.AnyArg(), // created_at
			).
			WillReturnResult(sqlmock.NewResult(int64(i), 1))
	}

	// Run the function
	err = PopulateWorkflows(db)

	// Assertions
	assert.NoError(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations.")
}
