package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.RecordRepoer = new(RecordRepo)

// mock.UserRepo
type RecordRepo struct {
	records map[string]*model.Record
}

// Return a new UserRepo
func NewRecordRepo() *RecordRepo {
	return &RecordRepo{
		records: make(map[string]*model.Record),
	}
}

func (rr *RecordRepo) Create(record *model.Record) error {
	if _, ok := rr.records[record.Title]; ok {
		return errorType.ErrRecordAlreadyExist{
			RecordTitle: record.Title,
			Message:     "cannot create record that already exists",
			Location:    "RecordRepo.Create()",
		}
	}

	rr.records[record.Title] = record
	return nil
}

func (rr *RecordRepo) Get(title string) (*model.Record, error) {
	record, ok := rr.records[title]
	if !ok {
		return nil, errorType.ErrRecordNotFound{
			RecordTitle: title,
			Message:     "cannot get record not found",
			Location:    "RecordRepo.Get()",
		}
	}
	return record, nil
}

func (rr *RecordRepo) GetAll() []*model.Record {
	records := make([]*model.Record, 0, len(rr.records))
	for _, record := range rr.records {
		records = append(records, record)
	}
	return records
}

func (rr *RecordRepo) Update(title string, record *model.Record) error {
	if _, ok := rr.records[title]; !ok {
		return errorType.ErrRecordNotFound{
			RecordTitle: title,
			Message:     "cannot update record not found",
			Location:    "RecordRepo.Update()",
		}
	}
	if title != record.Title {
		return errorType.ErrRecordNotValid{
			RecordTitle: title,
			Message:     "cannot update record when title mismatch",
			Location:    "RecordRepo.Update()",
		}
	}

	rr.records[title] = record
	return nil
}

func (rr *RecordRepo) Delete(title string) error {
	if _, ok := rr.records[title]; !ok {
		return errorType.ErrRecordNotFound{
			RecordTitle: title,
			Message:     "cannot delete record not found",
			Location:    "RecordRepo.Delete()",
		}
	}
	delete(rr.records, title)
	return nil
}
