package main

import (
	"errors"
)

type CommandLineParser struct {
	Environment string
	Source      string
}

func (cpl *CommandLineParser) Parse(args []string) error {

	var err error
	cpl.Environment, err = cpl.GetEnvironment(args)
	if err != nil {
		return err
	}

	return nil
}

func (cpl CommandLineParser) GetEnvironment(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Please specify environment.")
	}
	return args[0], nil
}
