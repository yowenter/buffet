package spider

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yowenter/buffet/pkg/lib"
)

var CREATED = "Created"
var RUNNING = "Running"
var SUCCESS = "Success"
var Failure = "Failure"

var DOWNLOAD = "Download"
var PARSE = "Parse"
var DUMP = "Dump"

var TotalTasks = NewLookTasks(100)

type Task struct {
	Id        string      `json:"id"`
	Url       string      `json:"url"`
	State     string      `json:"state"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Phase     string      `json:"phase"`
	Logs      []string    `json:"logs"`
	Result    interface{} `json:"result"`
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

type TaskResponse struct {
	Id       string
	Response *http.Response
}

type TaskItem struct {
	Id   string
	Item *lib.Item
}

type TaskRequest struct {
	Id      string
	Request *http.Request
}

type TaskMsg struct {
	Id      string
	Phase   string
	Message string
	Data    interface{}
}

type LoopTasks struct {
	Tasks []*Task `json:"tasks"`
	Idx   int
	Size  int
}

func NewLookTasks(size int) *LoopTasks {
	tasks := make([]*Task, size)
	loopTasks := LoopTasks{
		Tasks: tasks,
		Idx:   0,
		Size:  size,
	}
	return &loopTasks
}

func (lt *LoopTasks) PushTask(task *Task) {
	if lt.Idx < lt.Size {
		lt.Tasks[lt.Idx] = task
		lt.Idx++
	} else {
		lt.Idx = 0
		lt.Tasks[lt.Idx] = task
		lt.Idx++
	}
}

func (lt *LoopTasks) UpdateTask(taskMsg *TaskMsg) {

	for _, task := range lt.Tasks {
		if task == nil {
			return
		}
		if task.Id == taskMsg.Id {
			task.Phase = taskMsg.Phase
			task.Logs = append(task.Logs, taskMsg.Message)
			if task.Phase == DOWNLOAD || task.Phase == PARSE {
				task.State = RUNNING
			} else if task.Phase == DUMP {
				task.State = SUCCESS
				task.Result = taskMsg.Data
			}
		}
	}
}
