package test

import (
	"github.com/golang/glog"
	"gopkg.in/gcfg.v1"
)

type ConfigAuto struct {
	Commands struct {
		Reads     int
		Conflicts int
	}
}

func ParseAuto(filename string) ConfigAuto {
	var config ConfigAuto
	err := gcfg.ReadFileInto(&config, filename)
	if err != nil {
		glog.Fatalf("Failed to parse gcfg data: %s", err)
	}
	return config
}