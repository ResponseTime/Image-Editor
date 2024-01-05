package Models

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
	TimeCreated string
}

type Resizer struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}
