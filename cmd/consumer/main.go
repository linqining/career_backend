package main

import (
	"encoding/json"
	"fmt"
	"os"

	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
	logger "github.com/go-eagle/eagle/pkg/log"
	v "github.com/go-eagle/eagle/pkg/version"
	"github.com/spf13/pflag"

	"career_backend/internal/tasks"
)

var (
	cfgDir  = pflag.StringP("config dir", "c", "config", "config path.")
	env     = pflag.StringP("env name", "e", "", "env var name.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()
	if *version {
		ver := v.Get()
		marshaled, err := json.MarshalIndent(&ver, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

	// init config
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	logger.Init(logger.WithFilename("consumer"))

	// load config
	c = config.New(*cfgDir, config.WithEnv(*env))
	var taskCfg tasks.Config
	if err := c.Load("consumer", &taskCfg); err != nil {
		panic(err)
	}

	// start app
	app, cleanup, err := InitApp(&cfg, &cfg.HTTP, &taskCfg)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		panic(err)
	}
}
