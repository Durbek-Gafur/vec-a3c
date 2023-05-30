package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"vec-node/internal/api"
	rs "vec-node/internal/rspec"
	rspec_mock "vec-node/internal/rspec/mocks"
	"vec-node/internal/store"
	store_mock "vec-node/internal/store/mocks"
)

func TestHandler_GetQueueSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueSizeStore(ctrl)
	h := api.NewHandler(mockStore, nil,nil)
	ctx := context.TODO()
	
	req, err := http.NewRequest("GET", "/queue-size", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := 5
	mockStore.EXPECT().GetQueueSize(ctx).Return(expectedSize,nil)

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

	var response map[string]float64
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


/*
func TestHandler_SetQueueSize(t *testing.T) {
	mockQueueStore := new(mocks.QueueSizeStore)
	h := api.NewHandler(mockQueueStore, nil, nil)

	size := map[string]string{"size": "10"}
	sizeJSON, _ := json.Marshal(size)

	req, err := http.NewRequest("POST", "/queue-size", bytes.NewBuffer(sizeJSON))
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := 10
	mockQueueStore.On("SetQueueSize", mock.Anything, expectedSize).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SetQueueSize)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedSize, response["size"])
	mockQueueStore.AssertCalled(t, "SetQueueSize", mock.Anything, expectedSize)
}

func TestHandler_UpdateQueueSize(t *testing.T) {
	mockQueueStore := new(mocks.QueueSizeStore)
	h := api.NewHandler(mockQueueStore, nil, nil)

	size := map[string]string{"size": "20"}
	sizeJSON, _ := json.Marshal(size)

	req, err := http.NewRequest("PUT", "/queue-size", bytes.NewBuffer(sizeJSON))
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := 20
	mockQueueStore.On("UpdateQueueSize", mock.Anything, expectedSize).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateQueueSize)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, strconv.Itoa(expectedSize), response["size"])
	mockQueueStore.AssertCalled(t, "UpdateQueueSize", mock.Anything, expectedSize)
}

// Tests for WorkflowStore

func TestHandler_GetWorkflowByID(t *testing.T) {
	mockWorkflowStore := new(mocks.WorkflowStore)
	h := api.NewHandler(nil, mockWorkflowStore, nil)

	req, err := http.NewRequest("GET", "/workflow/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedWorkflow := store.Workflow{ID: 1, Name: "Workflow1"}
	mockWorkflowStore.On("GetWorkflowByID", mock.Anything, expectedWorkflow.ID).Return(&expectedWorkflow, nil)

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/workflow/{id}", h.GetWorkflowByID)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response store.Workflow
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedWorkflow, response)
	mockWorkflowStore.AssertCalled(t, "GetWorkflowByID", mock.Anything, expectedWorkflow.ID)
}

func TestHandler_GetWorkflows(t *testing.T) {
	mockWorkflowStore := new(mocks.WorkflowStore)
	h := api.NewHandler(nil, mockWorkflowStore, nil)

	req, err := http.NewRequest("GET", "/workflows", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedWorkflows := []store.Workflow{
		{ID: 1, Name: "Workflow1"},
		{ID: 2, Name: "Workflow2"},
	}
	mockWorkflowStore.On("GetWorkflows", mock.Anything, mock.Anything).Return(expectedWorkflows, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetWorkflows)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []store.Workflow
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedWorkflows, response)
	mockWorkflowStore.AssertCalled(t, "GetWorkflows", mock.Anything, mock.Anything)
}

func TestHandler_SaveWorkflow(t *testing.T) {
	mockWorkflowStore := new(mocks.WorkflowStore)
	h := api.NewHandler(nil, mockWorkflowStore, nil)

	wf := store.Workflow{ID: 1, Name: "Workflow1"}
	wfJSON, _ := json.Marshal(wf)

	req, err := http.NewRequest("POST", "/workflow", bytes.NewBuffer(wfJSON))
	if err != nil {
		t.Fatal(err)
	}

	mockWorkflowStore.On("SaveWorkflow", mock.Anything, &wf).Return(&wf, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SaveWorkflow)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response store.Workflow
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, wf, response)
	mockWorkflowStore.AssertCalled(t, "SaveWorkflow", mock.Anything, &wf)
}

// Test for Rspec

func TestHandler_GetRspec(t *testing.T) {
	mockRspec := new(mocks.Rspec)
	h := api.NewHandler(nil, nil, mockRspec)

	req, err := http.NewRequest("GET", "/rspec", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedRspec := store.Rspec{RAM: 1024, CPUs: 4}
	mockRspec.On("GetRspec").Return(&expectedRspec, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetRspec)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]float64
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, float64(expectedRspec.RAM), response["RAM"])
	assert.Equal(t, float64(expectedRspec.CPUs), response["CORE"])
	mockRspec.AssertCalled(t, "GetRspec")
}
*/