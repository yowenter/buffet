package spider

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yowenter/buffet/pkg/lib"
	"github.com/yowenter/buffet/pkg/plugins"
	"github.com/yowenter/buffet/pkg/plugins/storage/airtable"
)

// Spider can receive tasks from api server
// The downloader will download request and yield responses to spider
// the main components maybe like `scrapy` https://doc.scrapy.org/en/latest/topics/architecture.html
// but there're some components removed for simple
// 爬虫 task 有 download, parse, save 三个阶段， 每个阶段 通过 channel 更新 task 状态。
// Task 列表 存放在 内存中， 使用 hash 表 存储。  当 task 超过 100 个时， 将老的清除。

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

type Spider struct {
	TaskChan     chan *Task
	responseChan chan *TaskResponse
	itemPipeChan chan *TaskItem
	storage      plugins.Storage
	name         string
	taskMsgChan  chan *TaskMsg
	// todo 增加重试机制；
}

func (s *Spider) String() string {
	return fmt.Sprintf("Spider <%s>", s.name)
}

func NewSpider(n string) *Spider {
	taskCh := make(chan *Task, 10) // 避免阻塞的channel
	resCh := make(chan *TaskResponse, 10)
	itemPipeCh := make(chan *TaskItem, 10)
	taskMsgCh := make(chan *TaskMsg, 30)
	store := airtable.Airtable{}
	spider := Spider{
		TaskChan:     taskCh,
		responseChan: resCh,
		itemPipeChan: itemPipeCh,
		storage:      &store,
		name:         n,
		taskMsgChan:  taskMsgCh,
	}
	return &spider
}

func (spider *Spider) SendMsg(id string, phase string, msg string, data interface{}) {

	taskMsg := TaskMsg{
		Id:      id,
		Message: msg,
		Phase:   phase,
		Data:    data,
	}
	spider.taskMsgChan <- &taskMsg
	log.Debugf("New task message %+v", taskMsg)
}

func (spider *Spider) ManageTasks() {
	for true {
		taskMsg, ok := <-spider.taskMsgChan
		if !ok {
			log.Warnf("Get item from task message channel failure")
			time.Sleep(3 * time.Second)
			continue
		}
		log.Debugf("Manage task state by task message %+v", taskMsg)
	}
}

//Run  task consumer
func (spider *Spider) Run() {
	go func() {
		log.Info("Start consuming reuquest task")
		for true {
			task, ok := <-spider.TaskChan
			if !ok {
				log.Warnf("Fetch task from channel failed")
				time.Sleep(3 * time.Second)
				continue
			}
			log.Debugf("New task from api server %+v", *task)
			time.Sleep(100 * time.Millisecond)

			request, error := http.NewRequest("GET", task.Url, nil)
			if error != nil {
				log.Errorf("Build request error %v+", error)
			}
			taskRequest := TaskRequest{
				Id:      task.Id,
				Request: request,
			}
			go spider.download(&taskRequest)
		}
	}()

	go func() {
		log.Info("Start consuming response task")
		for true {
			response, ok := <-spider.responseChan
			if !ok {
				log.Warnf("Fetch response from channel failed ")
				time.Sleep(3 * time.Second)
				continue
			}

			item := spider.parse(response)

			if item == nil {
				continue
			}
			spider.itemPipeChan <- item
		}
	}()

	go func() {
		log.Info("Start consuming item pipe task")
		for true {
			item, ok := <-spider.itemPipeChan
			if !ok {
				log.Warnf("Fetch item pipe from channel failed")
				continue
			}
			go spider.dump(item)

		}

	}()

	go spider.ManageTasks()

}

func (spider *Spider) download(taskReq *TaskRequest) {
	taskReq.Request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(taskReq.Request)
	if err != nil {
		log.Errorf("Fetching request failure %v+ err:", taskReq.Request, err.Error())
		return
	}
	log.Debugf("Download request %v, %v", taskReq.Request.URL, res.Status)
	taskResp := TaskResponse{
		Id:       taskReq.Id,
		Response: res,
	}
	spider.responseChan <- &taskResp
	spider.SendMsg(taskReq.Id, "Download", fmt.Sprintf("Downloading %s", taskReq.Request.URL), nil)
}

func (spider *Spider) parse(taskResp *TaskResponse) *TaskItem {

	log.Debugf("Start parsing response from %v", taskResp.Response.Request.URL)

	// you can implement your own plugins
	// accoring to your website
	plugin := plugins.MatchPlugin(taskResp.Response.Request.URL)

	if plugin == nil {
		log.Infof("Url %v parser not found", taskResp.Response.Request.URL)
		return nil
	}

	item := (*plugin).Parse(taskResp.Response)
	taskItem := TaskItem{
		Id:   taskResp.Id,
		Item: item,
	}
	spider.SendMsg(taskResp.Id, "Parse", fmt.Sprintf("Parsed %s", taskResp.Id), taskItem)
	return &taskItem

}

func (spider *Spider) dump(taskItem *TaskItem) {
	// Dump item to airtable & Google Spreadsheet
	// Maybe
	log.Debugf("Dumping item %+v", taskItem.Item)
	spider.storage.Dump(taskItem.Item)
	spider.SendMsg(taskItem.Id, "Dump", fmt.Sprintf("Dumped %s", taskItem.Id), nil)
}

//
//
