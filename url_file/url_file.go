/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package url_file

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

func createDir(path string) error {
	if _, err1 := os.Stat(path); os.IsNotExist(err1) {
		if err2 := os.MkdirAll(path, 0777); err2 != nil {
			return err2
		}
	} else {
		return err1
	}

	return nil
}

func createURLFilePath(urlStr, rootPath string) string {
	filePath := url.QueryEscape(urlStr)
	filePath = path.Join(rootPath, filePath)

	return filePath
}

func SaveWebPage(rootPath, url string, data []byte) error {
	if err := createDir(rootPath); err != nil {
		return fmt.Errorf("createDir() err: %s", err.Error())
	}

	filePath := createURLFilePath(url, rootPath)

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("ioutil.WriteFile() err: %s", err.Error())
	}

	return nil
}
