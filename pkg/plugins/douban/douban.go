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

	// the selectors are copied from chrome -)
	// I don't know much about jquery .
	title := doc.Find("#wrapper > h1 > span")
	author := doc.Find("#info > span:nth-child(1) > a")
	intro := doc.Find("#link-report > span.all.hidden > div > div")
	var tags = []string{}

	tagsSelector := doc.Find("#db-tags-section > div")
	tagsSelector.Each(
		// #db-tags-section > div > span:nth-child(1) > a
		func(i int, s *goquery.Selection) {
			aSelector := s.Find("span >a")
			tags = append(tags, aSelector.Text())
		})

	return &lib.Item{
		Author:      author.Text(),
		Link:        *(res.Request.URL),
		Subject:     title.Text(),
		Description: intro.Text(),
		Tags:        tags,
	}

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
