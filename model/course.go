package model

import "time"

type Course struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	Subjects     map[int]string `json:"subjects"`
	Cover        string         `json:"cover"`
	Description  string         `json:"description"`
	Duration     string         `json:"duration"`
	TeacherId    uint           `json:"teacher_id"`
	ViewCount    uint           `json:"view_count"`
	ClassId      uint           `json:"class_id"`
	CourseDetail string         `json:"course_detail"`
	Date         time.Time      `json:"date"`
	FavoriteBy   []*User        `json:"favorite_by"`
}
