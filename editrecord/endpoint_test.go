package editrecord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"interviewtest/record"
	"interviewtest/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const tmpStorageFilePath = "/tmp/edit_records.bin"

func TestEditRecordSuccessful(t *testing.T) {
	t.Parallel()

	fileStorageService, err := storage.NewService(tmpStorageFilePath)

	if err != nil {
		t.Error(err)
	}

	defer os.Remove(tmpStorageFilePath)
	defer fileStorageService.Close()

	service := NewService(fileStorageService)

	testingTime := time.Date(2023, 12, 31, 12, 42, 59, 987654321, time.Local)

	record1 := record.Record{
		Id:        1,
		IntValue:  42,
		StrValue:  "foo1",
		BoolValue: false,
		TimeValue: &testingTime,
	}

	record2 := record.Record{
		Id:        2,
		IntValue:  90,
		StrValue:  "foo2",
		BoolValue: false,
		TimeValue: &testingTime,
	}

	expectedRecord1 := record.Record{
		Id:        1,
		IntValue:  142,
		StrValue:  "foo+1",
		BoolValue: true,
		TimeValue: &testingTime,
	}

	expectedRecord2 := record.Record{
		Id:        2,
		IntValue:  190,
		StrValue:  "foo+2",
		BoolValue: false,
		TimeValue: &testingTime,
	}

	expectedRecordId1, err := fileStorageService.CreateRecord(&record1)

	if err != nil {
		t.Fatal(err)
	}

	expectedRecordId2, err := fileStorageService.CreateRecord(&record2)

	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		service Service
	}

	tests := []struct {
		name         string
		args         args
		originalData record.Record
		expectedData record.Record
		idURLParam   int64
	}{{
		name:         "Edit record by id: 1",
		args:         args{service: service},
		originalData: record1,
		expectedData: expectedRecord1,
		idURLParam:   expectedRecordId1,
	}, {
		name:         "Edit record by id: 2",
		args:         args{service: service},
		originalData: record2,
		expectedData: expectedRecord2,
		idURLParam:   expectedRecordId2,
	}}

	for _, tt := range tests {
		url := fmt.Sprintf("/records/%d", tt.idURLParam)

		router := mux.NewRouter()
		router.Handle("/records/{id:[0-9]+}", MakePutRecordEndpoint(service)).Methods(http.MethodPut)

		inputDataBytes, _ := json.Marshal(tt.expectedData)

		req, _ := http.NewRequest("PUT", url, bytes.NewReader(inputDataBytes))

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		persistedData, err := fileStorageService.GetRecord(tt.idURLParam)

		if err != nil {
			t.Error(err)
		}

		assert.NotNil(t, persistedData)
		assert.NotEqual(t, int64(0), persistedData.Id)
		assert.Equal(t, tt.expectedData.Id, persistedData.Id)
		assert.Equal(t, tt.expectedData.IntValue, persistedData.IntValue)
		assert.Equal(t, tt.expectedData.StrValue, persistedData.StrValue)
		assert.Equal(t, tt.expectedData.BoolValue, persistedData.BoolValue)
		assert.Equal(t, tt.expectedData.TimeValue, persistedData.TimeValue)
	}
}

func TestEditRecordInternalServerError(t *testing.T) {
	fileStorageService, err := storage.NewService("/tmp/edit_records.bin")

	if err != nil {
		t.Error(err)
	}

	defer t.Cleanup(func() {
		os.Remove("/tmp/tmp_records.bin")
		fileStorageService.Close()
	})

	service := NewService(fileStorageService)

	testingTime := time.Date(2023, 12, 31, 12, 42, 59, 987654321, time.Local)

	record1 := record.Record{
		Id:        1,
		IntValue:  42,
		StrValue:  "foo1",
		BoolValue: false,
		TimeValue: &testingTime,
	}

	fileStorageService.CreateRecord(&record1)

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
		inputData  string
		idURLParam int64
	}{{
		name:       "Edit record - expected status code 500",
		args:       args{service: service},
		inputData:  rec,
		idURLParam: int64(1),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			router.Handle("/records/{id:[0-9]+}", MakePutRecordEndpoint(service)).Methods(http.MethodPut)

			req, _ := http.NewRequest("PUT", fmt.Sprintf("/records/%d", tt.idURLParam), strings.NewReader(tt.inputData))

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusInternalServerError, rr.Code)
		})
	}
}
