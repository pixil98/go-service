package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type cmdLineOpts struct {
	config    string
	loglevel  string
	logformat string

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

	if opts.fs.Lookup("logformat") != nil {
		return nil, fmt.Errorf("logformat flag already set")
	}
	opts.fs.StringVar(&opts.logformat, "logformat", "json", "logger output format (json, text)")

	return opts, nil
}

func (o *cmdLineOpts) Parse(args []string) error {
	err := o.fs.Parse(args[1:])
	if err != nil {
		return err
	}

	return nil
}

func (o *cmdLineOpts) Logger() (*slog.Logger, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(o.loglevel))
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", o.loglevel, err)
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch o.logformat {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stderr, opts)
	default:
		return nil, fmt.Errorf("invalid log format %q: must be json or text", o.logformat)
	}

	return slog.New(handler), nil
}

func (o *cmdLineOpts) Config(cfg any) error {
	if o.config == "" {
		return fmt.Errorf("config flag is required")
	}

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
