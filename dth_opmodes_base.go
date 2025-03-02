package main

import (
	"github.com/sirupsen/logrus"
)

type DthAutoToolOpMode interface {
	Run() error
}

type DthAutoToolOpModeBase struct {
	logger *logrus.Logger
	cfg    DthAutoToolConfig
}

func (o *DthAutoToolOpModeBase) SetLogger(logger *logrus.Logger) {
	o.logger = logger
}

func (o *DthAutoToolOpModeBase) SetConfiguration(cfg DthAutoToolConfig) {
	o.cfg = cfg
}
