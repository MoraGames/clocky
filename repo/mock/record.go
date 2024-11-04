package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// RecordRepo Error
type ErrRecordRepo struct {
	RecordId int64
	Message  string
	Location string
}

func (err ErrRecordRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.RecordId)
}

// Check if the repo implements the interface
var _ repo.RecordRepoer = new(RecordRepo)

// RecordRepo is a mock implementation
type RecordRepo struct {
	records map[int64]*model.Record
	lastId  int64
}

func NewRecordRepo() *RecordRepo {
	return &RecordRepo{
		records: make(map[int64]*model.Record),
		lastId:  -1,
	}
}

func (rr *RecordRepo) Create(record *model.Record) (int64, error) {
	if rr.titleAlreadyUsed(record) {
		return -1, ErrRecordRepo{
			RecordId: -1,
			Message:  "record title already used",
			Location: "RecordRepo.Create()",
		}
	}

	rr.lastId++
	record.ID = rr.lastId
	rr.records[rr.lastId] = record
	return rr.lastId, nil
}

func (rr *RecordRepo) Get(id int64) (*model.Record, error) {
	record, ok := rr.records[id]
	if !ok {
		return nil, ErrRecordRepo{
			RecordId: id,
			Message:  "record not found",
			Location: "RecordRepo.Get()",
		}
	}
	return record, nil
}

func (rr *RecordRepo) GetAll() []*model.Record {
	bonuses := make([]*model.Record, 0, len(rr.records))
	for _, bonus := range rr.records {
		bonuses = append(bonuses, bonus)
	}
	return bonuses
}

func (rr *RecordRepo) Update(id int64, record *model.Record) error {
	_, ok := rr.records[id]
	if !ok {
		return ErrRecordRepo{
			RecordId: id,
			Message:  "record not found",
			Location: "RecordRepo.Update()",
		}
	}
	if id != record.ID {
		return ErrRecordRepo{
			RecordId: id,
			Message:  "records id mismatch",
			Location: "RecordRepo.Update()",
		}
	}

	if rr.titleAlreadyUsed(record) {
		return ErrRecordRepo{
			RecordId: id,
			Message:  "record title already used",
			Location: "RecordRepo.Update()",
		}
	}

	rr.records[id] = record
	return nil
}

func (rr *RecordRepo) Delete(id int64) error {
	_, ok := rr.records[id]
	if !ok {
		return ErrRecordRepo{
			RecordId: id,
			Message:  "record not found",
			Location: "RecordRepo.Delete()",
		}
	}
	delete(rr.records, id)
	return nil
}

func (rr *RecordRepo) titleAlreadyUsed(record *model.Record) bool {
	for _, e := range rr.records {
		if e.Title == record.Title {
			return true
		}
	}
	return false
}
