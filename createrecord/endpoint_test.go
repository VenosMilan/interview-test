package createrecord

import (
	"bytes"
	"encoding/json"
	"interviewtest/record"
	"interviewtest/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const tmpStorageFilePath = "/tmp/create_records.bin"

func TestCreationValidRecords(t *testing.T) {
	fileStorageService, err := storage.NewService(tmpStorageFilePath)

	if err != nil {
		t.Error(err)
	}

	defer t.Cleanup(func() {
		os.Remove(tmpStorageFilePath)
		fileStorageService.Close()
	})

	service := NewService(fileStorageService)

	testingTime := time.Date(2023, 12, 31, 12, 42, 59, 987654321, time.Local)

	record1 := record.Record{
		IntValue:  42,
		StrValue:  "foo",
		BoolValue: false,
		TimeValue: &testingTime,
	}

	record2 := record.Record{
		IntValue:  99,
		StrValue:  "foo99",
		BoolValue: true,
		TimeValue: &testingTime,
	}

	record3 := record.Record{
		IntValue:  150,
		StrValue:  "foo150",
		BoolValue: true,
		TimeValue: &testingTime,
	}

	type args struct {
		service Service
	}

	tests := []struct {
		name       string
		args       args
		returnedID bool
		inputData  record.Record
	}{{
		name:       "Create record with ID 1",
		args:       args{service: service},
		returnedID: true,
		inputData:  record1,
	}, {
		name:       "Create record with ID 2",
		args:       args{service: service},
		returnedID: true,
		inputData:  record2,
	}, {
		name:       "Create record with ID 3",
		args:       args{service: service},
		returnedID: true,
		inputData:  record3,
	}}

	for _, tst := range tests {
		tt := tst
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inputDataBytes, _ := json.Marshal(tt.inputData)

			req, err := http.NewRequest("POST", "/create", bytes.NewReader(inputDataBytes))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := MakePostCreateRecordEndpoint(service)
			handler(rr, req)

			assert.Equal(t, rr.Code, http.StatusCreated)
			assert.Equal(t, rr.Header().Get("Content-Type"), "application/json")

			var response map[string]int64
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.returnedID, response["ID"] != 0)
		})
	}
}

func TestCreateRecordInternalServerError(t *testing.T) {
	fileStorageService, err := storage.NewService(tmpStorageFilePath)

	if err != nil {
		t.Error(err)
	}

	defer t.Cleanup(func() {
		os.Remove(tmpStorageFilePath)
		fileStorageService.Close()
	})

	service := NewService(fileStorageService)

	rec := `{
		"IntValue": 42,
		"StrValue": "Test Record",
		"BoolValue": false
	}`

	type args struct {
		service Service
	}

	tests := []struct {
		name       string
		args       args
		returnedID bool
		inputData  string
	}{{
		name:       "Expected status code 500",
		args:       args{service: service},
		returnedID: true,
		inputData:  rec,
	}}

	for _, tst := range tests {
		tt := tst
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/create", strings.NewReader(tt.inputData))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := MakePostCreateRecordEndpoint(service)
			handler(rr, req)

			assert.Equal(t, http.StatusInternalServerError, rr.Code)
		})
	}
}

func TestCreationNonValidRecordBadRequest(t *testing.T) {
	fileStorageService, err := storage.NewService(tmpStorageFilePath)

	if err != nil {
		t.Error(err)
	}

	defer t.Cleanup(func() {
		os.Remove(tmpStorageFilePath)
		fileStorageService.Close()
	})

	service := NewService(fileStorageService)

	rec := `{
		"IntValue": 42,
		"StrValue": 44,
		"BoolValue": false,
		"TimeValue": "2023-10-14T12:00:00Z"
	}`

	type args struct {
		service Service
	}

	tests := []struct {
		name       string
		args       args
		returnedID bool
		inputData  string
	}{{
		name:       "Expected status code 400",
		args:       args{service: service},
		returnedID: true,
		inputData:  rec,
	}}

	for _, tst := range tests {
		tt := tst
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/create", strings.NewReader(tt.inputData))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := MakePostCreateRecordEndpoint(service)
			handler(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
	}
}
