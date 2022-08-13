/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package mini_spider

import (
	"github.com/sirupsen/logrus"

	"github.com/Hs88888/mini-spider/config_load"
	"github.com/Hs88888/mini-spider/crawler"
	"github.com/Hs88888/mini-spider/queue"
	"github.com/Hs88888/mini-spider/url_table"
)

type Spider struct {
	config   *config_load.Config
	urlTable *url_table.URLTable
	Que      *queue.Queue
	crawlers []*crawler.Crawler
}

func NewSpider(urls []string, config *config_load.Config) *Spider {
	spider := new(Spider)
	spider.config = config
	spider.urlTable = url_table.NewURLTable()
	spider.Que = new(queue.Queue)
	spider.Que.Init()
	spider.crawlers = make([]*crawler.Crawler, 0)

	var err error
	for _, url := range urls {
		err = spider.urlTable.Add(url)
		if err != nil {
			logrus.WithField("err", err.Error()).Info("spider.urlTable.Add failed")
			continue
		}

		crawlTask := &crawler.CrawlTask{
			URL:   url,
			Depth: 0,
		}
		err = spider.Que.Append(crawlTask)
		if err != nil {
			logrus.WithField("err", err.Error()).Error("spider.Que.Append failed")
			continue
		}
	}

	for i := 0; i < config.Basic.ThreadCnt; i++ {
		crawler, err := crawler.NewCrawler(spider.urlTable, spider.config, spider.Que)
		if err != nil {
			logrus.WithField("err", err.Error()).Error("crawler.NewCrawler() failed")
			continue
		}
		spider.crawlers = append(spider.crawlers, crawler)
	}

	return spider
}

func (spider *Spider) Run() {
	for _, crawler := range spider.crawlers {
		go crawler.Run()
	}
	logrus.Info("spider.Run()")
}

func (spider *Spider) Wait() {
	spider.Que.Join()
}
