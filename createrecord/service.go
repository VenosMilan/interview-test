package createrecord

import (
	"interviewtest/record"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// Service interface provides method for creating record in file storage
type Service interface {
	Create(rec *record.Record) (*createResponse, error)
}

type service struct {
	record record.ModificationStorage
}

// NewService constructor of service
// Argument is interface of storage
func NewService(record record.Storage) Service {
	return &service{record: record}
}

// Create method creating validating and creating record
// Method return response with id of new record
func (service *service) Create(rec *record.Record) (*createResponse, error) {
	validate := validator.New()

	if err := validate.Struct(rec); err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	id, err := service.record.CreateRecord(rec)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &createResponse{RecordID: id}, nil
}

type createResponse struct {
	RecordID int64 `json:"ID"`
}
