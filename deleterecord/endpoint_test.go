package deleterecord

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"interviewtest/record"
	"interviewtest/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const tmpStorageFilePath = "/tmp/delete_records.bin"

func TestDeleteRecordSuccessful(t *testing.T) {
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

	expectedRecord1 := record.Record{
		Id:        0,
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
		name:         "Delete record by id: 1",
		args:         args{service: service},
		originalData: record1,
		expectedData: expectedRecord1,
		idURLParam:   int64(1),
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
			router.Handle("/records/{id:[0-9]+}", MakeDeleteRecordEndpoint(service)).Methods(http.MethodDelete)

			req, _ := http.NewRequest("DELETE", url, nil)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusNoContent, rr.Code)

			persistedData := readFile(int64(0), file)

			assert.Equal(t, tt.expectedData.Id, persistedData.Id)
			assert.Equal(t, tt.expectedData.IntValue, persistedData.IntValue)
			assert.Equal(t, tt.expectedData.StrValue, persistedData.StrValue)
			assert.Equal(t, tt.expectedData.BoolValue, persistedData.BoolValue)
			assert.Equal(t, tt.expectedData.TimeValue, persistedData.TimeValue)
		})
	}
}

func TestDeleteRecordNotFound(t *testing.T) {
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

	file, err := os.OpenFile(tmpStorageFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		t.Error(err)
	}

	writerRecord(&record1, file)

	file.Close()

	type args struct {
		service Service
	}

	tests := []struct {
		name       string
		args       args
		idURLParam int64
	}{{
		name:       "Delete record by id: 99",
		args:       args{service: service},
		idURLParam: int64(99),
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
			router.Handle("/records/{id:[0-9]+}", MakeDeleteRecordEndpoint(service)).Methods(http.MethodDelete)

			req, _ := http.NewRequest("DELETE", url, nil)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusNotFound, rr.Code)
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
