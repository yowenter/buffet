package spider

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yowenter/buffet/pkg/lib"
	"github.com/yowenter/buffet/pkg/plugins"
)

type Task struct {
	Url  string
	Tags string
}

//Spider can receive tasks from api server
// The downloader will download request and yield responses to spider
// the main components maybe like `scrapy` https://doc.scrapy.org/en/latest/topics/architecture.html
// but there're some components removed for simple

type Spider struct {
	TaskChan     chan *Task
	responseChan chan *http.Response
}

func NewSpider() *Spider {
	taskCh := make(chan *Task, 10) // 避免阻塞的channel
	resCh := make(chan *http.Response, 10)
	spider := Spider{
		TaskChan:     taskCh,
		responseChan: resCh,
	}
	return &spider
}

//Run  task consumer
func (spider *Spider) Run() {
	go func() {
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

			go spider.download(request)
		}
	}()

	go func() {
		for true {
			response, ok := <-spider.responseChan
			if !ok {
				log.Warnf("Fetch response from channel failed ")
				continue
			}

			item := spider.parse(response)
			if item == nil {
				continue
			}
			spider.dump(item)

		}
	}()

}

func (spider *Spider) download(request *http.Request) {
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(request)
	if err != nil {
		log.Errorf("Fetching request failure %v+ err:", request, err.Error())
		return
	}
	log.Debugf("Download request %v, %v", request.URL, res.Status)
	spider.responseChan <- res
	// 阻塞调用.
	// 不同的网站下载方式是否有所不同？
	// suppport plugins

}

func (spider *Spider) parse(response *http.Response) *lib.Item {

	log.Debugf("Start parsing response from %v", response.Request.URL)

	// you can implement your own plugins
	// accoring to your website
	plugin := plugins.MatchPlugin(response.Request.URL)

	if plugin == nil {
		log.Infof("Url %v parser not found", response.Request.URL)
		return nil
	}

	item := (*plugin).Parse(response)

	return item

}

func (spider *Spider) dump(item *lib.Item) {
	// Dump item to airtable & Google Spreadsheet
	//
	log.Debugf("Dumping item %+v", item)

}

//
//
