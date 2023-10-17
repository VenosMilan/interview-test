package getrecord

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"interviewtest/record"
	"interviewtest/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulGetRecords(t *testing.T) {
	tmpStorageFilePath := "/tmp/tmp_get_records.bin"

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
		Id:        1,
		IntValue:  42,
		StrValue:  "foo",
		BoolValue: false,
		TimeValue: &testingTime,
	}

	record2 := record.Record{
		Id:        2,
		IntValue:  90,
		StrValue:  "foo2",
		BoolValue: true,
		TimeValue: &testingTime,
	}

	file, err := os.OpenFile(tmpStorageFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		t.Error(err)
	}

	writerRecord(&record1, file)
	writerRecord(&record2, file)

	type args struct {
		service Service
	}

	tests := []struct {
		name       string
		args       args
		expected   record.Record
		idURLParam int64
	}{{
		name:       "Read record by id: 1",
		args:       args{service: service},
		expected:   record1,
		idURLParam: int64(1),
	}, {
		name:       "Read record by id: 2",
		args:       args{service: service},
		expected:   record2,
		idURLParam: int64(2),
	}}

	for _, tst := range tests {
		tt := tst
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			url := fmt.Sprintf("/records/%d", tt.idURLParam)

			router := mux.NewRouter()
			router.Handle("/records/{id:[0-9]+}", MakeGetRecordEndpoint(service)).Methods(http.MethodGet)

			req, _ := http.NewRequest("GET", url, nil)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			var responseRecord record.Record

			if err = json.NewDecoder(rr.Body).Decode(&responseRecord); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			assert.NotEqual(t, int64(0), responseRecord.Id)
			assert.Equal(t, tt.expected.Id, responseRecord.Id)
			assert.Equal(t, tt.expected.IntValue, responseRecord.IntValue)
			assert.Equal(t, tt.expected.StrValue, responseRecord.StrValue)
			assert.Equal(t, tt.expected.BoolValue, responseRecord.BoolValue)
			assert.Equal(t, tt.expected.TimeValue, responseRecord.TimeValue)
		})
	}
}

func TestGetRecordNotFound(t *testing.T) {
	tmpStorageFilePath := "/tmp/tmp_get_records.bin"

	fileStorageService, err := storage.NewService(tmpStorageFilePath)

	if err != nil {
		t.Error(err)
	}

	defer t.Cleanup(func() {
		os.Remove(tmpStorageFilePath)
		fileStorageService.Close()
	})

	service := NewService(fileStorageService)

	type args struct {
		service Service
	}

	tests := []struct {
		name       string
		args       args
		idURLParam int64
	}{{
		name:       "Not found by id: 99",
		args:       args{service: service},
		idURLParam: 99,
	}}

	for _, tst := range tests {
		tt := tst
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			router.Handle("/records/{id:[0-9]+}", MakeGetRecordEndpoint(service)).Methods(http.MethodGet)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/records/%d", tt.idURLParam), nil)
			req2, _ := http.NewRequest("GET", fmt.Sprintf("/records/%s", "xyz"), nil)

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusNotFound, rr.Code)

			router.ServeHTTP(rr, req2)
			assert.Equal(t, http.StatusNotFound, rr.Code)
		})
	}
}

// writeRecord function for test
func writerRecord(rec *record.Record, file *os.File) error {
	binary.Write(file, binary.LittleEndian, rec.Id)
	binary.Write(file, binary.LittleEndian, rec.IntValue)

	strBytes := []byte(rec.StrValue)
	strBytes = append(strBytes, make([]byte, 64-len(strBytes))...)
	file.Write(strBytes)

	boolByte := byte(0)
	if rec.BoolValue {
		boolByte = byte(1)
	}

	file.Write([]byte{boolByte})

	timeBytes, _ := rec.TimeValue.MarshalBinary()
	timeLen := len(timeBytes)

	if timeLen < 16 {
		padding := make([]byte, 16-timeLen)
		timeBytes = append(timeBytes, padding...)
	}

	file.Write(timeBytes)
	file.Write([]byte{'\n'})

	return nil
}
