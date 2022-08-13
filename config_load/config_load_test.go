/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package config_load

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {
	filePath := "../config_file/config.conf"
	config, err := ConfigLoad(filePath)
	assert.Nil(t, err)
	assert.NotEqual(t, nil, &config)
}
