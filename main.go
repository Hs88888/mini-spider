/**
  @author: hs
  @date: 2022/8/13
  @note:

modification history
--------------------
**/

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Hs88888/mini-spider/config_load"
	"github.com/Hs88888/mini-spider/mini_spider"
	"github.com/Hs88888/mini-spider/seed_file"
)

var (
	help     *bool   = flag.Bool("h", false, "help")
	showVer  *bool   = flag.Bool("v", false, "version")
	confPath *string = flag.String("c", "./config_file/config.conf", "mini-spider config file path")
	logPath  *string = flag.String("l", "./log/mini-spider.log", "dir path of log")
	logLevel *int    = flag.Int("d", 3, "level of log")
)

var version = "dev"

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	if *showVer {
		fmt.Printf("mini-spider: version %s\n", version)
		return
	}

	logInit()

	config, err := config_load.ConfigLoad(*confPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err.Error()}).Error("load config failed")
		Exit(1)
	}
	logrus.Info("load config success")

	urls, err := seed_file.LoadSeedFile(config.Basic.URLFilePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err.Error()}).Error("load seed file failed")
		Exit(1)
	}
	logrus.Info("load seed file success")

	spider := mini_spider.NewSpider(urls, &config)
	spider.Run()

	fmt.Println("**************** mini-spider start ****************")
	spider.Wait()
	fmt.Println("**************** mini-spider done ****************")
	Exit(0)
}

func Exit(exitCode int) {
	time.Sleep(100 * time.Millisecond)
	os.Exit(exitCode)
}

func logInit() {
	switch *logLevel {
	case 1:
		logrus.SetLevel(logrus.TraceLevel)
	case 2:
		logrus.SetLevel(logrus.DebugLevel)
	case 3:
		logrus.SetLevel(logrus.InfoLevel)
	case 4:
		logrus.SetLevel(logrus.WarnLevel)
	case 5:
		logrus.SetLevel(logrus.ErrorLevel)
	case 6:
		logrus.SetLevel(logrus.FatalLevel)
	case 7:
		logrus.SetLevel(logrus.PanicLevel)
	default:
		Exit(1)
	}

	file, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("os.OpenFile() err: %s\n", err.Error())
		Exit(1)
	}
	logrus.SetOutput(file)
}
