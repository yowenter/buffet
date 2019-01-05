package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/yowenter/buffet/pkg/server"
	"github.com/yowenter/buffet/pkg/spider"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	fmt.Println("IFTTT Service Buffet Starting ....")
	iftttServiceKey := os.Getenv("IFTTT_SERVICE_KEY")
	iftttChannelKey := os.Getenv("IFTTT_CHANNEL_KEY")

	iftttConf := &server.IftttConf{
		IftttChannelKey: iftttChannelKey,
		IftttServiceKey: iftttServiceKey,
	}

	router := mux.NewRouter()
	spider := spider.NewSpider("Ifttt")
	buffetServer := &server.BuffetAPIServer{
		Router:    router,
		Spider:    spider,
		IftttConf: iftttConf,
	}
	log.Debugf("Init Buffet server:  `%+v` ", *buffetServer)
	buffetServer.InstallHandlers()

	buffetAPIServer := &http.Server{
		Addr:           ":5000",
		Handler:        buffetServer,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go buffetServer.Spider.Run()

	log.Fatal(buffetAPIServer.ListenAndServe())

}
