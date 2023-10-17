package getrecord

import (
	"interviewtest/record"
	"interviewtest/tools"

	"github.com/pkg/errors"
)

// Service interface provides method for reading record from file storage
type Service interface {
	GetRecord(id int64) (*record.Record, error)
}

type service struct {
	record record.ReadingStorage
}

// NewService constructor of service
// Argument is interface of storage
func NewService(record record.Storage) Service {
	return &service{record: record}
}

// GetRecord method for read record by id
func (service *service) GetRecord(id int64) (*record.Record, error) {
	rec, err := service.record.GetRecord(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if rec == nil {
		return nil, tools.RecordNotFound
	}

	return rec, nil
}
