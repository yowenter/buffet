package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/yowenter/buffet/pkg/spider"
)

func (s *BuffetAPIServer) home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Buffet server  %v", *s)
}

func (s *BuffetAPIServer) collect(w http.ResponseWriter, r *http.Request) {
	ok, errString := s.verifyIftttKey(r)
	if !ok {
		http.Error(w, errString, 401)
		return
	}

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
	s.Spider.TaskChan <- &task

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
