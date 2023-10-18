package storage

import (
	"interviewtest/record"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRecord(t *testing.T) {
	t.Parallel()

	tmpfile, err := os.CreateTemp("", "get_record.bin")
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

	rec := record.Record{
		IntValue:  9,
		StrValue:  "Test 1",
		BoolValue: true,
		TimeValue: &testingTime,
	}

	createdID, err := service.CreateRecord(&rec)
	assert.NoError(t, err)
	assert.NotZero(t, createdID)

	record1, err := service.GetRecord(createdID)
	assert.NoError(t, err)
	assert.NotNil(t, record1)
	assert.Equal(t, int64(1), record1.Id)
	assert.Equal(t, "Test 1", record1.StrValue)
	assert.Equal(t, true, record1.BoolValue)
	assert.Equal(t, testingTime, *record1.TimeValue)
	assert.Equal(t, int64(9), record1.IntValue)
}

func TestCreateRecords(t *testing.T) {
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
		assert.NotZero(t, createdID)

		lastID, err := service.getLastIdOfRecord(tmpfile)

		assert.NoError(t, err)
		assert.Equal(t, lastID, createdID)
		assert.Equal(t, newRecord.Id, lastID)
	}
}

func TestDeleteRecord(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "delete_records.bin")
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
		assert.NotZero(t, createdID)

		deleted, err := service.DeleteRecord(createdID)

		assert.NoError(t, err)
		assert.Equal(t, true, deleted)
	}
}

func TestGetRecordNotFound(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "get_record.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	service := &service{
		storageFilePath: tmpfile.Name(),
		storageFile:     tmpfile,
	}

	recordNotFound, err := service.GetRecord(42)
	assert.NoError(t, err)
	assert.Nil(t, recordNotFound)
}

func TestEditRecord(t *testing.T) {
	t.Parallel()

	tmpfile, err := os.CreateTemp("", "edit_records.bin")
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

	rec1 := record.Record{
		Id:        1,
		IntValue:  42,
		StrValue:  "foo",
		BoolValue: false,
		TimeValue: &testingTime,
	}

	createdID, err := service.CreateRecord(&rec1)
	assert.NoError(t, err)
	assert.NotZero(t, createdID)

	testingTimeForUpdate := time.Date(2022, 12, 31, 12, 42, 59, 987654321, time.Local)

	update1 := record.Record{
		IntValue:  142,
		StrValue:  "bee",
		BoolValue: true,
		TimeValue: &testingTimeForUpdate,
	}

	updateID, err := service.EditRecord(createdID, &update1)

	assert.NoError(t, err)
	assert.NotZero(t, updateID)

	updatedRecord, err := service.GetRecord(createdID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedRecord)
	assert.Equal(t, updateID, updatedRecord.Id)
	assert.Equal(t, "bee", updatedRecord.StrValue)
	assert.Equal(t, true, updatedRecord.BoolValue)
	assert.Equal(t, testingTimeForUpdate, *updatedRecord.TimeValue)
	assert.Equal(t, int64(142), updatedRecord.IntValue)
}
