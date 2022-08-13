/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Hs88888/mini-spider/config_load"
	"github.com/Hs88888/mini-spider/html_parse"
	"github.com/Hs88888/mini-spider/queue"
	"github.com/Hs88888/mini-spider/url_file"
	"github.com/Hs88888/mini-spider/url_table"
)

type Crawler struct {
	urlTable   *url_table.URLTable
	config     *config_load.ConfBasic
	que        *queue.Queue
	urlPattern *regexp.Regexp
}

type CrawlTask struct {
	URL   string
	Depth int
}

func NewCrawler(urlTable *url_table.URLTable, config *config_load.Config, queue *queue.Queue) (*Crawler, error) {
	crawler := new(Crawler)
	crawler.urlTable = urlTable
	crawler.config = &config.Basic
	crawler.que = queue

	var err error
	crawler.urlPattern, err = regexp.Compile(config.Basic.TargetURLReg)
	if err != nil {
		return crawler, fmt.Errorf("regexp.Compile() err: %s", err.Error())
	}

	return crawler, nil
}

func (crawler *Crawler) Run() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("err", err).Panic("crawler.Run()")
		}
	}()

	for {
		crawlTask := crawler.que.Remove().(*CrawlTask)

		data, err := getPageData(crawlTask.URL)
		if err != nil {
			logrus.WithField("err", err.Error()).Warn("getPageData() err")
			crawler.que.TaskDone()
			continue
		}

		if crawler.urlPattern.MatchString(crawlTask.URL) {
			if err = url_file.SaveWebPage(crawler.config.OutputPath, crawlTask.URL, data); err != nil {
				logrus.WithField("err", err.Error()).Error("url_file.SaveWebPage() failed")
			}
		}

		if crawlTask.Depth < crawler.config.MaxDepth {
			err = crawler.addChildTasks(data, crawlTask)
			if err != nil {
				logrus.WithField("err", err.Error()).Error("crawler.addChildTasks failed")
			}
		}

		crawler.que.TaskDone()
		time.Sleep(time.Duration(crawler.config.Interval) * time.Second)
	}
}

func getPageData(url string) ([]byte, error) {
	var data []byte
	var err error

	resp, err := http.Get(url)
	if err != nil {
		return data, fmt.Errorf("http.Get failed, err: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, fmt.Errorf("ioutil.ReadAll failed, err: %s", err.Error())
	}

	return data, nil
}

func (crawler *Crawler) addChildTasks(data []byte, crawlTask *CrawlTask) error {
	links, err := html_parse.ParseWebPage(crawlTask.URL, data)
	if err != nil {
		return fmt.Errorf("html_parse.ParseWebPage() err: %s", err.Error())
	}

	for _, link := range links {
		if err := crawler.urlTable.Add(link); err != nil {
			logrus.WithField("err", err.Error()).Warn("crawler.urlTable.Add() err")
			continue
		}

		newCrawlTask := &CrawlTask{
			URL:   link,
			Depth: crawlTask.Depth + 1,
		}
		crawler.que.Append(newCrawlTask)
	}

	return nil
}
