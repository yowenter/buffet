package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/yowenter/buffet/pkg/spider"
)

// ifttt configuration
type IftttConf struct {
	IftttServiceKey string
	IftttChannelKey string
}

// BuffetAPIServer ...
type BuffetAPIServer struct {
	router    *mux.Router
	spider    *spider.Spider
	IftttConf *IftttConf
}

// type IftttExampleData struct {
//     Data struct {
//         AccessToken      string `json:"accessToken"`

//     } `json:"data"`
// }

// type IftttSamples struct {
// 	Actions struct {
// 		Collect struct {
// 			Url string `json:"url"`
// 			Tags string `json:"tags"`
// 		} `json:"collect"`

// 	} `json:"actions"`
// } `json:"samples"`

type Collect struct {
	URL  string `json:"url"`
	Tags string `json:"tags"`
}

type Actions struct {
	Collect Collect `json:"collect_a_link"`
}

type Samples struct {
	Actions Actions `json:"actions"`
}
type IftttSamplesData struct {
	Samples Samples `json:"samples"`
}

type IftttTestData struct {
	Data IftttSamplesData `json:"data"`
}

type IftttMessage struct {
	ActionFields Collect `json:"actionFields"`
}

func main() {
	fmt.Println("Ifttt service `buffet` starting ....")
	iftttServiceKey := os.Getenv("IFTTT_SERVICE_KEY")
	iftttChannelKey := os.Getenv("IFTTT_CHANNEL_KEY")

	iftttConf := &IftttConf{
		IftttChannelKey: iftttChannelKey,
		IftttServiceKey: iftttServiceKey,
	}

	router := mux.NewRouter()
	spider := spider.NewSpider()
	buffetServer := &BuffetAPIServer{
		router:    router,
		spider:    spider,
		IftttConf: iftttConf,
	}
	log.Debugf("Init Buffet server  %+v", *buffetServer)
	buffetServer.InstallHandlers()

	buffetAPIServer := &http.Server{
		Addr:           ":5000",
		Handler:        buffetServer,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go buffetServer.spider.Run()

	log.Fatal(buffetAPIServer.ListenAndServe())

}

func (s *BuffetAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.router.ServeHTTP(w, r)
	end := time.Now()
	log.WithFields(log.Fields{"Method": r.Method, "Path": r.URL, "Addr": r.RemoteAddr, "Elapsed": end.Sub(start)}).Info("")
}

func (s *BuffetAPIServer) InstallHandlers() {
	s.router.HandleFunc("/", s.home)
	s.router.HandleFunc("/ifttt/v1/actions/collect_a_link", s.collect)
	s.router.HandleFunc("/ifttt/v1/test/setup", s.test)
	s.router.HandleFunc("/ifttt/v1/status", s.status)
}

func (s *BuffetAPIServer) home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Buffet server  %v", *s)
}

func (s *BuffetAPIServer) collect(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var iftttMsg IftttMessage
	err = json.Unmarshal(b, &iftttMsg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Debugf("IftttMsg URL: %v, Tags: %v ", iftttMsg.ActionFields.URL, iftttMsg.ActionFields.Tags)

	task := spider.Task{
		Url:  iftttMsg.ActionFields.URL,
		Tags: iftttMsg.ActionFields.Tags,
	}
	s.spider.TaskChan <- &task

	response := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "success",
		Message: "",
	}
	b, er := json.Marshal(response)
	if er != nil {
		http.Error(w, er.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(b))
}

func (s *BuffetAPIServer) verifyIftttKey(r *http.Request) (b bool, errString string) {
	verified := true
	errString = ""
	serviceKey, ok := r.Header["Ifttt-Service-Key"]
	if !ok {
		errString = "no service key"
		return false, errString
	}
	channelKey, ok := r.Header["Ifttt-Channel-Key"]
	if !ok {
		errString = "no channel key"
		return false, errString
	}

	if strings.Join(serviceKey, "") != s.IftttConf.IftttServiceKey || strings.Join(channelKey, "") != s.IftttConf.IftttChannelKey {
		errString = "invalid service or channel key"
	}
	if len(errString) > 0 {
		verified = false
	}
	return verified, errString
}

func (s *BuffetAPIServer) status(w http.ResponseWriter, r *http.Request) {
	ok, err := s.verifyIftttKey(r)
	if !ok {
		http.Error(w, err, 401)
		return
	}
	fmt.Fprintf(w, "ok")
}

func (s *BuffetAPIServer) test(w http.ResponseWriter, r *http.Request) {
	ok, errString := s.verifyIftttKey(r)
	if !ok {
		http.Error(w, errString, 401)
		return
	}

	// those codes below is so stupid :-)
	exampleCollect := Collect{
		URL:  "http://blog.heytaoge.com",
		Tags: "Blog, Personal",
	}
	exampleActions := Actions{
		Collect: exampleCollect,
	}
	exampleSamples := Samples{
		Actions: exampleActions,
	}
	exampleData := IftttSamplesData{
		Samples: exampleSamples,
	}
	data := IftttTestData{
		Data: exampleData,
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}
	fmt.Fprintf(w, string(b))

}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}
