package Models

import "time"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserStack struct {
	UndoStack    []*Image
	RedoStack    []*Image
	CurrentImage *Image
}
type Image struct {
	Path string
}

type Result struct {
	ImagePath   string
	User        string
	ProjectName string
	Date        string
	TimeCreated time.Time
}
