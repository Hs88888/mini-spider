/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package config_load

import (
	"fmt"

	"github.com/go-gcfg/gcfg"
)

type ConfBasic struct {
	URLFilePath  string
	OutputPath   string
	MaxDepth     int
	Interval     int
	Timeout      int
	ThreadCnt    int
	TargetURLReg string
}

type Config struct {
	Basic ConfBasic
}

func ConfigLoad(filePath string) (Config, error) {
	var cfg Config

	if err := gcfg.ReadFileInto(&cfg, filePath); err != nil {
		return cfg, fmt.Errorf("gcfg.ReadFileInto err: %s", err.Error())
	}

	return cfg, nil
}

