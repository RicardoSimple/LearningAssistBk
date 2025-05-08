package model

import "time"

type AssignmentSubmission struct {
	Id           uint      `json:"id"`
	AssignmentId uint      `json:"assignment_id"`
	StudentId    uint      `json:"student_id"`
	Content      string    `json:"content"`
	Title        string    `json:"title"`
	Score        float64   `json:"score"`
	FeedBack     string    `json:"feedback"`
	SubmittedAt  time.Time `json:"submitted_at"`
	ReviewedAt   time.Time `json:"reviewed_at"`
}
