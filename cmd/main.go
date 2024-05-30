package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"twitch-chatbot/internal/configurations"
	"twitch-chatbot/internal/plugins"
	"twitch-chatbot/internal/twitch"
	"twitch-chatbot/internal/ui"
)

var (
	banner string = `  _____          _ _       _        ____ _           _   ____        _
 |_   _|_      _(_) |_ ___| |__    / ___| |__   __ _| |_| __ )  ___ | |_
   | | \ \ /\ / / | __/ __| '_ \  | |   | '_ \ / _  | __|  _ \ / _ \| __|
   | |  \ V  V /| | || (__| | | | | |___| | | | (_| | |_| |_) | (_) | |_
   |_|   \_/\_/ |_|\__\___|_| |_|  \____|_| |_|\__,_|\__|____/ \___/ \__|
`
	titleView bool
)

func main() {
	var configsPath string
	// Executable Flags
	flag.StringVar(&configsPath, "c", "configs.json", "Path to configuration json file")
	flag.BoolVar(&titleView, "t", false, "Flags to keep start banner always visable")
	flag.Parse()

	log.Println("Loading configurations")
	configs, err := configurations.Load(configsPath)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Printf("%+v\n", configs)

	log.Println("Loading messages filters")
	filters := plugins.LoadFilter(configs.Filters)

	log.Println("Loading Commands macros")
	commands := plugins.LoadCommands(configs.Commands)

	// create ui configs
	mainUI := ui.NewUI(titleView, banner)
	uiStarted := make(chan struct{}, 1)

	// create twitch client
	client := twitch.NewClient(
		configs.TwitchIRL,
		configs.Channel,
		configs.Debug,
		mainUI.WriteMain,
		mainUI.WriteCmd,
		mainUI.WriteSide,
		commands,
		filters,
	)
	defer client.Close()

	// start twitch client
	go client.StartBot(uiStarted)

	// config watchdog
	go func() {
		oldStat, _ := os.Stat(configsPath)
		for {
			stat, _ := os.Stat(configsPath)
			if oldStat.ModTime() != stat.ModTime() {
				// reload file
				if changed := configs.Reload(configsPath); changed {
					// reload twitch client
					client.ReloadConfig(
						uiStarted,
						configs.TwitchIRL,
						configs.Channel,
						configs.Debug,
					)
					// reload plugins
					client.Commands = plugins.LoadCommands(configs.Commands)
					client.Filters = plugins.LoadFilter(configs.Filters)
					// update prints
					client.WriteCurrentConfigs()
				}
				oldStat = stat
			}
			time.Sleep(time.Second)
		}
	}()

	// start graphica interface
	mainUI.Start(uiStarted)
}
