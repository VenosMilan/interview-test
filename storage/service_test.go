package storage

import (
	"encoding/binary"
	"interviewtest/record"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRecord(t *testing.T) {
	t.Parallel()

	tmpfile, err := os.CreateTemp("", "get_records.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testingTime := time.Date(2023, 12, 31, 12, 42, 59, 987654321, time.Local)

	records := []record.Record{
		{Id: 1, IntValue: 42, StrValue: "Test 1", BoolValue: true, TimeValue: &testingTime},
		{Id: 2, IntValue: 99, StrValue: "Test 2", BoolValue: false, TimeValue: &testingTime},
	}

	for _, r := range records {
		if err := binary.Write(tmpfile, binary.LittleEndian, r.Id); err != nil {
			t.Fatal(err)
		}
		if err := binary.Write(tmpfile, binary.LittleEndian, r.IntValue); err != nil {
			t.Fatal(err)
		}
		strBytes := []byte(r.StrValue)
		strBytes = append(strBytes, make([]byte, 64-len(strBytes))...)
		if _, err := tmpfile.Write(strBytes); err != nil {
			t.Fatal(err)
		}
		boolByte := []byte{0}
		if r.BoolValue {
			boolByte[0] = 1
		}
		if _, err := tmpfile.Write(boolByte); err != nil {
			t.Fatal(err)
		}
		timeBytes, err := r.TimeValue.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		timeBytes = append(timeBytes, make([]byte, 16-len(timeBytes))...)
		if _, err := tmpfile.Write(timeBytes); err != nil {
			t.Fatal(err)
		}
		if _, err := tmpfile.Write([]byte{0x0A}); err != nil {
			t.Fatal(err)
		}
	}

	service := &service{
		storageFilePath: tmpfile.Name(),
		storageFile:     tmpfile,
	}

	record1, err := service.GetRecord(1)
	assert.NoError(t, err)
	assert.NotNil(t, record1)
	assert.Equal(t, int64(1), record1.Id)
	assert.Equal(t, "Test 1", record1.StrValue)
	assert.Equal(t, true, record1.BoolValue)
	assert.Equal(t, testingTime, *record1.TimeValue)
	assert.Equal(t, int64(42), record1.IntValue)

	record2, err := service.GetRecord(2)
	assert.NoError(t, err)
	assert.NotNil(t, record2)
	assert.Equal(t, int64(2), record2.Id)
	assert.Equal(t, "Test 2", record2.StrValue)
	assert.Equal(t, false, record2.BoolValue)
	assert.Equal(t, testingTime, *record2.TimeValue)
	assert.Equal(t, int64(99), record2.IntValue)

	recordNotFound, err := service.GetRecord(42)
	assert.NoError(t, err)
	assert.Nil(t, recordNotFound)
}

func TestCreateRecord(t *testing.T) {
	t.Parallel()

	tmpfile, err := os.CreateTemp("", "create_records.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	service := &service{
		storageFilePath: tmpfile.Name(),
		storageFile:     tmpfile,
	}

	testingTime := time.Date(2023, 12, 31, 12, 42, 59, 987654321, time.Local)

	testingRecords := make([]record.Record, 0)
	testingRecords = append(testingRecords, record.Record{
		IntValue:  42,
		StrValue:  "foo",
		BoolValue: false,
		TimeValue: &testingTime,
	})

	testingRecords = append(testingRecords, record.Record{
		IntValue:  99,
		StrValue:  "foo99",
		BoolValue: true,
		TimeValue: &testingTime,
	})

	for _, newRecord := range testingRecords {
		createdID, err := service.CreateRecord(&newRecord)
		assert.NoError(t, err)
		assert.NotNil(t, createdID)

		lastID, err := service.getLastIdOfRecord(tmpfile)

		assert.NoError(t, err)
		assert.Equal(t, lastID, createdID)
		assert.Equal(t, newRecord.Id, lastID)
	}
}
