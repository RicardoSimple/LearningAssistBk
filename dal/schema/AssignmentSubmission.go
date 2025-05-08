package schema

import (
	"learning-assistant/model"
	"time"
)

type AssignmentSubmission struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	AssignmentID uint       `gorm:"not null" json:"assignment_id"` // 作业ID
	StudentID    uint       `gorm:"not null" json:"student_id"`    // 提交人（学生）
	Content      string     `gorm:"type:text" json:"content"`      // 提交的内容（支持 Markdown / 纯文本 / 链接）
	Title        string     `gorm:"type:text" json:"title"`
	Score        float64    `json:"score"`                     // 成绩（教师填写）
	Feedback     string     `gorm:"type:text" json:"feedback"` // 批改反馈
	SubmittedAt  time.Time  `json:"submitted_at"`              // 提交时间
	ReviewedAt   *time.Time `json:"reviewed_at"`               // 教师批改时间
}

func (as *AssignmentSubmission) ToType() *model.AssignmentSubmission {

	asm := &model.AssignmentSubmission{
		Id:           as.ID,
		AssignmentId: as.AssignmentID,
		StudentId:    as.StudentID,
		Content:      as.Content,
		Title:        as.Title,
		Score:        as.Score,
		FeedBack:     as.Feedback,
		SubmittedAt:  as.SubmittedAt,
	}
	if as.ReviewedAt != nil {
		asm.ReviewedAt = *as.ReviewedAt
	}
	return asm
}
