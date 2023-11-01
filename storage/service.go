package storage

import (
	"bytes"
	"encoding/binary"
	"interviewtest/record"
	"io"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const recordSize = 98

// Service interface that provides method for working with binary file storage
type Service interface {
	GetRecord(id int64) (*record.Record, error)
	CreateRecord(rec *record.Record) (int64, error)
	EditRecord(id int64, rec *record.Record) (int64, error)
	DeleteRecord(id int64) (bool, error)
	Close()
}

type service struct {
	storageFilePath string
	storageFile     *os.File
	mu              sync.Mutex
}

// NewService constructor for create new binary file storage
// constructor create binary file
func NewService(fileStoragePath string) (Service, error) {
	file, err := os.OpenFile(fileStoragePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &service{
		storageFilePath: fileStoragePath,
		storageFile:     file,
	}, nil
}

// GetRecord method for get record by id from binary file
func (service *service) GetRecord(id int64) (*record.Record, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	recPos := (id - 1) * recordSize

	_, err := service.storageFile.Seek(recPos, io.SeekStart)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var rec record.Record

	if err := binary.Read(service.storageFile, binary.LittleEndian, &rec.Id); err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := binary.Read(service.storageFile, binary.LittleEndian, &rec.IntValue); err != nil {
		return nil, errors.WithStack(err)
	}

	strBytes := make([]byte, 64)
	if _, err := service.storageFile.Read(strBytes); err != nil {
		return nil, errors.WithStack(err)
	}

	rec.StrValue = string(bytes.TrimRight(strBytes, string(rune(0))))

	boolByte := make([]byte, 1)
	if _, err := service.storageFile.Read(boolByte); err != nil {
		return nil, errors.WithStack(err)
	}

	rec.BoolValue = boolByte[0] != 0

	timeBytes := make([]byte, 16)
	if _, err := service.storageFile.Read(timeBytes); err != nil {
		return nil, errors.WithStack(err)
	}

	timeBytes = bytes.Trim(timeBytes, "\x00")

	var t time.Time

	if err := t.UnmarshalBinary(timeBytes); err != nil {
		return nil, errors.WithStack(err)
	}

	rec.TimeValue = &t

	if _, err := service.storageFile.Seek(1, io.SeekCurrent); err != nil {
		return nil, errors.WithStack(err)
	}

	return &rec, nil
}

// CreateRecord method for create record in binary file
// ID of record is the defined position of the record in the file
// Every new file is appended to end of file
func (service *service) CreateRecord(rec *record.Record) (int64, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	pos, err := service.storageFile.Seek(0, io.SeekEnd)

	if err != nil {
		return 0, errors.WithStack(err)
	}

	rec.Id = pos/recordSize + 1

	if err := service.writeRecord(rec, service.storageFile); err != nil {
		return 0, errors.WithStack(err)
	}

	return rec.Id, nil
}

// EditRecord method for edit record in binary file by id and return updated id
func (service *service) EditRecord(id int64, updatedRecord *record.Record) (int64, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	recPos := (id - 1) * recordSize

	if _, err := service.storageFile.Seek(recPos, io.SeekStart); err != nil {
		return 0, errors.WithStack(err)
	}

	if err := service.writeRecord(updatedRecord, service.storageFile); err != nil {
		return 0, errors.WithStack(err)
	}

	return updatedRecord.Id, nil
}

// DeleteRecord method delete record by id (set id to zero)
// other data of record other data are not changed
func (service *service) DeleteRecord(id int64) (bool, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	_, err := service.storageFile.Seek(0, io.SeekStart)
	if err != nil {
		return false, errors.WithStack(err)
	}

	var deletedRecord bool

	for {
		var actualId int64

		if err := binary.Read(service.storageFile, binary.LittleEndian, &actualId); err == io.EOF {
			break
		} else if err != nil {
			return false, errors.WithStack(err)
		}

		if actualId == id {
			// offset -8 for rewrite id of record
			_, err = service.storageFile.Seek(-8, io.SeekCurrent)

			if err != nil {
				return false, errors.WithStack(err)
			}

			if err := binary.Write(service.storageFile, binary.LittleEndian, int64(0)); err != nil {
				return false, errors.WithStack(err)
			}

			deletedRecord = true
		} else {
			// offset 90 - skip to next record in file
			_, err = service.storageFile.Seek(90, io.SeekCurrent)

			if err != nil {
				return false, errors.WithStack(err)
			}
		}
	}

	return deletedRecord, nil
}

// Close method close storage file
func (service *service) Close() {
	service.storageFile.Close()
}

func (service *service) writeRecord(rec *record.Record, file *os.File) error {
	if err := binary.Write(file, binary.LittleEndian, rec.Id); err != nil {
		return errors.WithStack(err)
	}

	if err := binary.Write(file, binary.LittleEndian, rec.IntValue); err != nil {
		return errors.WithStack(err)
	}

	strBytes := make([]byte, 64)
	copy(strBytes, rec.StrValue)

	if _, err := file.Write(strBytes); err != nil {
		return errors.WithStack(err)
	}

	boolByte := byte(0)
	if rec.BoolValue {
		boolByte = byte(1)
	}

	if _, err := file.Write([]byte{boolByte}); err != nil {
		return errors.WithStack(err)
	}

	timeBytes, err := rec.TimeValue.MarshalBinary()

	if err != nil {
		return errors.WithStack(err)
	}

	timeLen := len(timeBytes)

	if timeLen < 16 {
		padding := make([]byte, 16-timeLen)
		timeBytes = append(timeBytes, padding...)
	}

	if _, err := file.Write(timeBytes); err != nil {
		return errors.WithStack(err)
	}

	if _, err := file.Write([]byte{'\n'}); err != nil {
		return errors.WithStack(err)
	}

	log.Debugf("Record %+v", rec)
	return nil
}
