package model

import "time"

type Course struct {
	ID          uint
	Name        string
	Subjects    map[int]string
	Cover       string
	Description string
	Duration    string
	TeacherId   uint
	ClassId     uint
	Date        time.Time
}
