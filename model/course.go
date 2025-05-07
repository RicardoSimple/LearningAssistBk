package model

import "time"

type Course struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Subjects    map[int]string `json:"subjects"`
	Cover       string         `json:"cover"`
	Description string         `json:"description"`
	Duration    string         `json:"duration"`
	TeacherId   uint           `json:"teacher_id"`
	ClassId     uint           `json:"class_id"`
	Date        time.Time      `json:"date"`
}
