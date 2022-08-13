/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package seed_file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func LoadSeedFile(filePath string) ([]string, error) {
	urlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile err: %s", err.Error())
	}

	var urls []string
	err = json.Unmarshal(urlData, &urls)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal err: %s", err.Error())
	}

	return urls, nil
}
