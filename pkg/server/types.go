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
	Collect Collect `json:"collect"`
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

type IftttAction struct {
	ActionFields Collect     `json:"actionFields"`
	IftttSource  IftttSource `json:"ifttt_source"`
	User         User        `json:"user"`
}

type User struct {
	Timezone string `json:"timezone"`
}

type IftttSource struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type IftttResp struct {
	Data []IftttObject `json:"data"`
}

type IftttObject struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Error struct {
	Message string `json:"message"`
}
type ErrorResp struct {
	Errors []Error `json:"errors"`
}
