package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/yowenter/buffet/pkg/spider"
)

// BuffetAPIServer ...
type BuffetAPIServer struct {
	airtableKey string
	airtableAPI string
	router      *mux.Router
	spider      *spider.Spider
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

type IftttMessage struct {
	ActionFields Collect `json:"actionFields"`
}

func main() {
	fmt.Println("Main")

	airtableKey := flag.String("airtable-key", "", "Airtable API Key")
	airtableAPI := flag.String("airtable-api", "", "Airtable base API url")
	flag.Parse()

	router := mux.NewRouter()
	spider := spider.NewSpider()
	buffetServer := &BuffetAPIServer{
		airtableKey: *airtableKey,
		airtableAPI: *airtableAPI,
		router:      router,
		spider:      spider,
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

func (s *BuffetAPIServer) test(w http.ResponseWriter, r *http.Request) {

	if val, ok := r.Header["Ifttt-Service-Key"]; ok {
		log.Debug("Ifttt-Service-key", val)

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

		b, err := json.Marshal(exampleData)
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
		fmt.Fprintf(w, string(b))

	} else {
		fmt.Fprintf(w, "No ifttt service key")
	}

}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}
