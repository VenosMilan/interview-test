package record

import "time"

// ReadingStorage interface provides methods for reading operations
type ReadingStorage interface {
	GetRecord(id int64) (*Record, error)
}

// ModificationStorage interface provides methods for modification operation
type ModificationStorage interface {
	CreateRecord(rec *Record) (int64, error)
	EditRecord(id int64, updatedRecord *Record) (int64, error)
	DeleteRecord(id int64) (bool, error)
}

// Storage interface provides access for reading and modification operation
type Storage interface {
	ReadingStorage
	ModificationStorage
}

type Record struct {
	Id        int64      `json:"id"`
	IntValue  int64      `json:"IntValue" validate:"required"`
	StrValue  string     `json:"StrValue" validate:"required"`
	BoolValue bool       `json:"BoolValue"`
	TimeValue *time.Time `json:"TimeValue" validate:"required"`
}
