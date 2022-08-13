/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package html_parse

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type HtmlLinks struct {
	links []string
}

func NewHtmlLinks() *HtmlLinks {
	hl := new(HtmlLinks)
	hl.links = make([]string, 0)

	return hl
}

func (hl *HtmlLinks) getLinks(n *html.Node, retURL *url.URL) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				linkURL, err := retURL.Parse(a.Val)
				if err != nil {
					logrus.WithField("err", err.Error()).Errorf("retURL.Parse() err")
				} else {
					hl.links = append(hl.links, linkURL.String())
				}
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		hl.getLinks(c, retURL)
	}
}

func ParseWebPage(urlStr string, data []byte) ([]string, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("html.Parse() err: %s", err.Error())
	}

	retURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, fmt.Errorf("url.ParseRequestURI() err: %s", err.Error())
	}

	hl := NewHtmlLinks()
	hl.getLinks(doc, retURL)

	return hl.links, nil
}