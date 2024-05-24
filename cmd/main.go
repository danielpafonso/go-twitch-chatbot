package main

import (
	"flag"
	"log"

	"twitch-chatbot/internal/configurations"
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
	// set file watchdog
	go configurations.FileWatch(configsPath, configs.Reload)

	// log.Println("Loading Commands macros")
	// plugins := LoadCommands()
	// client.ReadChat(plugins)

	// create ui configs
	ui := ui.NewUI(titleView, banner)
	uiStarted := make(chan struct{}, 1)

	// create twitch client
	client := twitch.NewClient(
		configs.TwitchIRL,
		configs.Channel,
		configs.Debug,
		ui.WriteMain,
		ui.WriteCmd,
		ui.WriteSide,
	)
	defer client.Close()

	// start twitch client
	go client.StartBot(uiStarted)

	// start graphica interface
	ui.Start(uiStarted)
}
