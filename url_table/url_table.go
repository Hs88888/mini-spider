/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package url_table

import (
	"fmt"
	"sync"
)

type URLTable struct {
	lock sync.Mutex
	table map[string]bool
}

func NewURLTable() *URLTable {
	urlTable := new(URLTable)
	urlTable.table = make(map[string]bool)

	return urlTable
}

func (urlTable *URLTable) Add(url string) error {
	urlTable.lock.Lock()
	defer urlTable.lock.Unlock()

	_, ok := urlTable.table[url]
	if ok {
		return fmt.Errorf("url[%s] exist", url)
	} else {
		urlTable.table[url] = true
		return nil
	}
}
