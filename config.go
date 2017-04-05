package main

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Root     string    `yaml:"root"`
	Listen   string    `yaml:"listen"`
	Flavors  []string  `yaml:"flavors"`
	Machines []Machine `yaml:"machines"`
}

func ReadConfig(path string) (Config, error) {
	var c Config

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return c, err
	}

	if c.Root == "" {
		c.Root = "/tftproot"
	}
	if c.Listen == "" {
		c.Listen = ":8085"
	}
	return c, nil
}

func (c Config) ValidFlavor(flavor string) bool {
	for _, test := range c.Flavors {
		if test == flavor {
			return true
		}
	}
	return false
}

func (c Config) ValidFlavors() string {
	return strings.Join(c.Flavors, ", ")
}
