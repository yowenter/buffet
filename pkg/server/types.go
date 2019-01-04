package server

import (
	"github.com/gorilla/mux"
	"github.com/yowenter/buffet/pkg/spider"
)

// ifttt configuration
type IftttConf struct {
	IftttServiceKey string
	IftttChannelKey string
}

// BuffetAPIServer ...
type BuffetAPIServer struct {
	Router    *mux.Router
	Spider    *spider.Spider
	IftttConf *IftttConf
}

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