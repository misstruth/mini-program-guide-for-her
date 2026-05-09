package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// StudyStore 学习数据存储接口
type StudyStore interface {
	ListTasks() ([]model.StudyTaskModel, error)
	CreateTask(task *model.StudyTaskModel) error
	UpdateTaskStatus(id uint, status string) error
	DeleteTask(id uint) error
	ListRecords(limit int) ([]model.StudyRecordModel, error)
	CreateRecord(record *model.StudyRecordModel) error
	UpdateRecord(record *model.StudyRecordModel) error
	DeleteRecord(id uint) error
	ListNotes(limit int) ([]model.StudyNoteModel, error)
	CreateNote(note *model.StudyNoteModel) error
	UpdateNote(note *model.StudyNoteModel) error
	DeleteNote(id uint) error
}

// GormStudyStore 基于 GORM 的学习数据存储
type GormStudyStore struct{}

// NewGormStudyStore 创建 GORM store
func NewGormStudyStore() StudyStore {
	return &GormStudyStore{}
}

func (s *GormStudyStore) ListTasks() ([]model.StudyTaskModel, error) {
	var tasks []model.StudyTaskModel
	err := db.Get().Order("sort asc, id asc").Find(&tasks).Error
	return tasks, err
}

func (s *GormStudyStore) CreateTask(task *model.StudyTaskModel) error {
	return db.Get().Create(task).Error
}

func (s *GormStudyStore) UpdateTaskStatus(id uint, status string) error {
	return db.Get().Model(&model.StudyTaskModel{}).Where("id = ?", id).Update("status", status).Error
}

func (s *GormStudyStore) DeleteTask(id uint) error {
	return db.Get().Delete(&model.StudyTaskModel{}, id).Error
}

func (s *GormStudyStore) ListRecords(limit int) ([]model.StudyRecordModel, error) {
	var records []model.StudyRecordModel
	query := db.Get().Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&records).Error
	return records, err
}

func (s *GormStudyStore) CreateRecord(record *model.StudyRecordModel) error {
	return db.Get().Create(record).Error
}

func (s *GormStudyStore) UpdateRecord(record *model.StudyRecordModel) error {
	return db.Get().Model(&model.StudyRecordModel{}).Where("id = ?", record.ID).Updates(map[string]interface{}{
		"title":    record.Title,
		"category": record.Category,
		"summary":  record.Summary,
	}).Error
}

func (s *GormStudyStore) DeleteRecord(id uint) error {
	return db.Get().Delete(&model.StudyRecordModel{}, id).Error
}

func (s *GormStudyStore) ListNotes(limit int) ([]model.StudyNoteModel, error) {
	var notes []model.StudyNoteModel
	query := db.Get().Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&notes).Error
	return notes, err
}

func (s *GormStudyStore) CreateNote(note *model.StudyNoteModel) error {
	return db.Get().Create(note).Error
}

func (s *GormStudyStore) UpdateNote(note *model.StudyNoteModel) error {
	return db.Get().Model(&model.StudyNoteModel{}).Where("id = ?", note.ID).Updates(map[string]interface{}{
		"title":   note.Title,
		"content": note.Content,
		"tag":     note.Tag,
	}).Error
}

func (s *GormStudyStore) DeleteNote(id uint) error {
	return db.Get().Delete(&model.StudyNoteModel{}, id).Error
}
