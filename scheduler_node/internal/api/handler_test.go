package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"scheduler-node/internal/api"
	"scheduler-node/internal/store"
	store_mock "scheduler-node/internal/store/mocks"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestShowTables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVenStore := store_mock.NewMockVENStore(ctrl)
	mockWorkflowStore := store_mock.NewMockWorkflowStore(ctrl)
	h := api.NewHandler(mockWorkflowStore, mockVenStore)
	ctx := context.TODO()

	// Set expectations on the mock
	mockVenStore.EXPECT().GetVENInfos().Return([]store.VENInfo{
		{
			Name:                    "Test VEN",
			RAM:                     "4GB",
			Core:                    "2",
			URL:                     "http://test.ven",
			MaxQueueSize:            "10",
			CurrentQueueSize:        "5",
			PreferenceList:          "Test Preference List",
			TrustScore:              "8",
			MaxQueueSizeUpdated:     time.Now(),
			CurrentQueueSizeUpdated: time.Now(),
			TrustScoreUpdated:       time.Now(),
		},
	}, nil)
	const layout = "2006-01-02 15:04:05"
	assignedAt, err := time.Parse(layout, "2023-05-31 10:00:00")
	if err != nil {
		panic(err)
	}
	assignedAt = assignedAt.UTC()

	completedAt, err := time.Parse(layout, "2023-05-31 10:00:00")
	if err != nil {
		panic(err)
	}
	completedAt = completedAt.UTC()
	mockWorkflowStore.EXPECT().GetWorkflows(ctx).Return([]store.WorkflowInfo{
		{
			Name:                  "Test WorkflowInfo",
			Type:                  "Sequential",
			RAM:                   "16GB",
			Core:                  "4",
			Policy:                "First-Come-First-Serve",
			ExpectedExecutionTime: "2 hours",
			ActualExecutionTime:   "1 hour",
			AssignedVM:            "VM1",
			SubmittedBy:           "TestUser",
			Status:                "pending",
			AssignedAt:            assignedAt,
			CompletedAt:           completedAt,
		},
	}, nil)

	req, err := http.NewRequest("GET", "/showtables", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ShowTables(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Due to the dynamic nature of the content it's complex to test the full HTML
	// Instead, we can check if the HTML contains certain parts
	expected := "Test VEN" // we are sure that the response should contain the name of the VEN
	assert.True(t,strings.Contains(rr.Body.String(), expected))

}

/*
func TestHandler_GetQueueSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.MockQueueStore(ctrl)
	h := api.NewHandler(mockStore, nil, nil)
	ctx := context.TODO()

	req, err := http.NewRequest("GET", "/queue-size", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := 5
	mockStore.EXPECT().GetQueueSize(ctx).Return(expectedSize, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetQueueSize)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedSize, response["size"])
}



func TestHandler_SetQueueSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueSizeStore(ctrl)
	h := api.NewHandler(mockStore, nil,nil)
	ctx := context.TODO()

	data := map[string]string{"size": "5"}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "/queue-size", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := 5
	mockStore.EXPECT().SetQueueSize(ctx, expectedSize).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SetQueueSize)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedSize, response["size"])
}

func TestHandler_SetQueueSize_InvalidSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueSizeStore(ctrl)
	h := api.NewHandler(mockStore, nil,nil)

	data := map[string]string{"size": "invalid"}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "/queue-size", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SetQueueSize)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}



func TestHandler_GetWorkflowByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockWorkflowStore(ctrl)
	h := api.NewHandler(nil, mockStore,nil)

	req, err := http.NewRequest("GET", "/workflow/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedWorkflow := &store.Workflow{ID: 1, Name: "Test workflow"}
	mockStore.EXPECT().GetWorkflowByID(gomock.Any(), 1).Return(expectedWorkflow, nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/workflow/{id}", h.GetWorkflowByID)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response store.Workflow
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedWorkflow, &response)
}

func TestHandler_UpdateQueueSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueSizeStore(ctrl)
	h := api.NewHandler(mockStore, nil,nil)
	ctx := context.TODO()

	data := map[string]string{"size": "5"}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "/queue-size", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := 5
	mockStore.EXPECT().UpdateQueueSize(ctx, expectedSize).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateQueueSize)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, strconv.Itoa(expectedSize), response["size"])
}

func TestHandler_SaveWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockWorkflowStore(ctrl)
	h := api.NewHandler(nil, mockStore,nil)
	ctx := context.TODO()

	wf := store.Workflow{ID: 1, Name: "TestWorkflow"}
	body, _ := json.Marshal(wf)
	req, err := http.NewRequest("POST", "/workflow", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	expectedWorkflow := wf
	mockStore.EXPECT().SaveWorkflow(ctx, &expectedWorkflow).Return(&expectedWorkflow, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SaveWorkflow)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response store.Workflow
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedWorkflow, response)
}

func TestHandler_GetRspec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRspec := rspec_mock.NewMockRspec(ctrl)
	h := api.NewHandler(nil, nil,mockRspec)

	req, err := http.NewRequest("GET", "/rspec", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedRspec := rs.Resources{RAM: "16", CPUs:"4"}
	mockRspec.EXPECT().GetRspec().Return(&expectedRspec, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetRspec)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedRspec.RAM, response["RAM"])
	assert.Equal(t, expectedRspec.CPUs, response["CORE"])
}


func TestHandler_GetWorkflows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockWorkflowStore(ctrl)
	h := api.NewHandler(nil, mockStore,nil)
	ctx := context.TODO()

	req, err := http.NewRequest("GET", "/workflows?type=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedWorkflows := []store.Workflow{{ID: 1, Name: "Workflow1"}, {ID: 2, Name: "Workflow2"}}
	filter := &store.WorkflowFilter{Type: "test"}
	mockStore.EXPECT().GetWorkflows(ctx, filter).Return(expectedWorkflows, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetWorkflows)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []store.Workflow
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedWorkflows, response)
}
*/