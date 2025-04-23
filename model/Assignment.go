package model

type Assignment struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	CourseId  uint    `json:"course_id"`
	Course    *Course `json:"course"`
	TeacherId uint    `json:"teacher_id"`
	Teacher   User    `json:"teacher"`
	DueDate   string  `json:"due_date"`
	CreatedAt string  `json:"created_at"`
}
