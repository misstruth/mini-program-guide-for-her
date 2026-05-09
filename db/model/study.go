package model

import "time"

// StudyTaskModel 学习任务
type StudyTaskModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(128);not null" json:"title"`
	Duration  string    `gorm:"type:varchar(32);not null" json:"duration"`
	Status    string    `gorm:"type:varchar(32);not null" json:"status"`
	Sort      int       `gorm:"not null;default:0" json:"sort"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// StudyRecordModel 学习记录
type StudyRecordModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(128);not null" json:"title"`
	Category  string    `gorm:"type:varchar(64);not null" json:"category"`
	Summary   string    `gorm:"type:text;not null" json:"summary"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// StudyNoteModel 复盘笔记
type StudyNoteModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(128);not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Tag       string    `gorm:"type:varchar(64);not null" json:"tag"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
