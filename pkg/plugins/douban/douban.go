package douban

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/yowenter/buffet/pkg/lib"
)

type DoubanParser struct {
	Name string
}

func (d *DoubanParser) Parse(res *http.Response) *lib.Item {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Warnf("Parse document error %v %v ", res, err)
	}
	//#info
	doc.Find("div#info").Each(func(i int, s *goquery.Selection) {
		log.Debug(i, s.Text())
	})

	log.Debug(doc.Find("div.related_info").Text())

	return &lib.Item{}

}

func (d *DoubanParser) Match(url *url.URL) bool {

	if url.Host == "book.douban.com" {

		matched, err := regexp.MatchString("/subject/(\\d+)/", url.Path)
		if err != nil {
			log.Warnf("Douban parser %v  match path %v error %v", d, url.Path, err.Error())
			return false
		}
		return matched

	}

	return false
}

func New() *DoubanParser {
	return &DoubanParser{Name: "DoubanReader Parser"}

}
