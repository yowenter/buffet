package plugins

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/yowenter/buffet/pkg/lib"
	"github.com/yowenter/buffet/pkg/plugins/douban"
)

type Parser interface {
	Parse(res *http.Response) *lib.Item
	Match(url *url.URL) bool
}

var plugins []Parser

func MatchPlugin(url *url.URL) *Parser {
	log.Debugf("Looking for plugin, host %v, path %v", url.Host, url.Path)
	for _, element := range plugins {

		if element.Match(url) {
			log.Debugf("Url %v matched plugin %v", url, element)
			return &element
		}
	}
	return nil
}

func registerPlugins() {
	doubanParser := douban.New()
	log.Infof("Register plugin douban %v", doubanParser)
	plugins = append(plugins, doubanParser)
}

func init() {
	registerPlugins()
}
