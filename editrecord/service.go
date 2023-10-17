package editrecord

import (
	"interviewtest/record"
	"interviewtest/tools"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// Service interface provides method for edit record in file storage
type Service interface {
	Edit(id int64, rec *record.Record) error
}

type service struct {
	record record.Storage
}

// NewService constructor of service
// Argument is interface of storage
func NewService(record record.Storage) Service {
	return &service{record: record}
}

// Edit method for validating and editing record
func (service *service) Edit(id int64, rec *record.Record) error {
	validate := validator.New()

	if err := validate.Struct(rec); err != nil {
		return err.(validator.ValidationErrors)
	}

	updateID, err := service.record.EditRecord(id, rec)

	if err != nil {
		return errors.WithStack(err)
	}

	if updateID == 0 {
		return tools.RecordNotFound
	}

	return nil
}
