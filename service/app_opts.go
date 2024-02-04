package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/pixil98/go-log/log"
)

type cmdLineOpts struct {
	config   string
	loglevel string

	fs flag.FlagSet
}

func newCmdLineOpts(name string) (*cmdLineOpts, error) {
	opts := &cmdLineOpts{}
	opts.fs = *flag.NewFlagSet(name, flag.ContinueOnError)

	if opts.fs.Lookup("config") != nil {
		return nil, fmt.Errorf("config flag already set")
	}
	opts.fs.StringVar(&opts.config, "config", "", "path to configuration file")

	if opts.fs.Lookup("loglevel") != nil {
		return nil, fmt.Errorf("loglevel flag already set")
	}
	opts.fs.StringVar(&opts.loglevel, "loglevel", "info", "logger log level (debug, info, warn, error)")

	return opts, nil
}

func (o *cmdLineOpts) Parse(args []string) error {
	err := o.fs.Parse(args[1:])
	if err != nil {
		return err
	}

	return nil
}

func (o *cmdLineOpts) Logger() (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(o.loglevel)
	if err != nil {
		return nil, err
	}

	var loggerOpts []log.LoggerOpt
	loggerOpts = append(loggerOpts, log.WithLevel(level))

	return log.NewLogger(loggerOpts...), nil
}

func (o *cmdLineOpts) Config(cfg interface{}) error {
	raw, err := os.ReadFile(o.config)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	err = json.Unmarshal(raw, cfg)
	if err != nil {
		return fmt.Errorf("unmarshaling config: %w", err)
	}
	return nil
}
