package main

import (
	"gitlab.com/cinemae/gopkg/log"
	"gopkg.in/yaml.v3"
)

func main() {

	fileConfig := []log.OutputConfig{
		{
			Writer:    "file",
			Formatter: "console",
			Level:     "error",
			WriteConfig: log.WriteConfig{
				LogPath:  "./",
				Filename: "test.log",
			},
			FormatConfig: log.FormatConfig{},
			RemoteConfig: yaml.Node{},
		},
		{
			Writer:    "console",
			Formatter: "json",
			Level:     "debug",
		},
	}

	fileZapLog := log.NewZapLog(fileConfig)
	log.Register("default", fileZapLog)
	defer fileZapLog.Sync()

	fileLog := log.GetLogger("default")

	fileLog.Info("test")
	fileLog.Errorf("test:%d", 100)
	fileLog.WithFields("name", "phachon").Info("with")
}
