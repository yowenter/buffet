package spider

import "time"
import "github.com/google/uuid"

var CREATED = "Created"
var RUNNING = "Running"
var SUCCESS = "Success"
var Failure = "Failure"

type Task struct {
	Id        string
	Url       string
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTask(url string) Task {
	id := uuid.New().String()
	createdAt := time.Now()
	updatedAt := time.Now()
	return Task{
		Id:        id,
		Url:       url,
		State:     CREATED,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
