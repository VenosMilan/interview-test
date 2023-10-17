package deleterecord

import (
	"interviewtest/record"
	"interviewtest/tools"

	"github.com/pkg/errors"
)

// Service interface provides method for delete record in file storage
type Service interface {
	Delete(id int64) error
}

type service struct {
	record record.Storage
}

// NewService constructor of service
// Argument is interface of storage
func NewService(record record.Storage) Service {
	return &service{record: record}
}

// Delete method for delete record by id
func (service *service) Delete(id int64) error {
	successfulDeleteRecord, err := service.record.DeleteRecord(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if !successfulDeleteRecord {
		return tools.RecordNotFound
	}

	return nil
}
