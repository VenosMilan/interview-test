package editrecord

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"interviewtest/record"
	"interviewtest/storage"
	"io"
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
		BoolValue: true,
		TimeValue: &testingTime,
	}

	file, err := os.OpenFile(tmpStorageFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		t.Error(err)
	}

	writerRecord(&record1, file)
	writerRecord(&record2, file)

	file.Close()

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
		idURLParam:   int64(1),
	}, {
		name:         "Edit record by id: 2",
		args:         args{service: service},
		originalData: record2,
		expectedData: expectedRecord2,
		idURLParam:   int64(2),
	}}

	file, err = os.Open(tmpStorageFilePath)

	if err != nil {
		t.Error(err)
	}

	for _, tst := range tests {
		tt := tst
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			url := fmt.Sprintf("/records/%d", tt.idURLParam)

			router := mux.NewRouter()
			router.Handle("/records/{id:[0-9]+}", MakePutRecordEndpoint(service)).Methods(http.MethodPut)

			inputDataBytes, _ := json.Marshal(tt.expectedData)

			req, _ := http.NewRequest("PUT", url, bytes.NewReader(inputDataBytes))

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			persistedData := readFile(tt.idURLParam, file)

			assert.NotEqual(t, int64(0), persistedData.Id)
			assert.Equal(t, tt.expectedData.Id, persistedData.Id)
			assert.Equal(t, tt.expectedData.IntValue, persistedData.IntValue)
			assert.Equal(t, tt.expectedData.StrValue, persistedData.StrValue)
			assert.Equal(t, tt.expectedData.BoolValue, persistedData.BoolValue)
			assert.Equal(t, tt.expectedData.TimeValue, persistedData.TimeValue)
		})
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

	file, err := os.OpenFile(tmpStorageFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		t.Error(err)
	}

	writerRecord(&record1, file)

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

// writeRecord function for test
func writerRecord(rec *record.Record, file *os.File) {
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

}

func readFile(targetID int64, file *os.File) *record.Record {
	_, _ = file.Seek(0, io.SeekStart)

	for {
		var rec record.Record

		if err := binary.Read(file, binary.LittleEndian, &rec.Id); err == io.EOF {
			return nil
		}

		binary.Read(file, binary.LittleEndian, &rec.IntValue)

		strBytes := make([]byte, 64)
		file.Read(strBytes)

		rec.StrValue = string(bytes.TrimRight(strBytes, string(rune(0))))

		boolByte := make([]byte, 1)
		file.Read(boolByte)

		rec.BoolValue = boolByte[0] != 0

		timeBytes := make([]byte, 16)
		file.Read(timeBytes)
		timeBytes = bytes.Trim(timeBytes, "\x00")

		var t time.Time

		t.UnmarshalBinary(timeBytes)

		rec.TimeValue = &t

		file.Seek(1, io.SeekCurrent)

		if rec.Id == targetID {
			return &rec
		}
	}
}
