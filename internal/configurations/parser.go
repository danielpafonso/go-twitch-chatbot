package configurations

import (
	"encoding/json"
	"os"
	"strings"

	"twitch-chatbot/internal/plugins"
)

type TwitchConfigs struct {
	TwitchIRL string                  `json:"twicthIrc"`
	Channel   string                  `json:"channel"`
	Debug     bool                    `json:"debug"`
	Filters   []plugins.FilterConfig  `json:"filter"`
	Commands  []plugins.CommandConfig `json:"commands"`
}

// Load reads and parse configuration json file
func Load(filepath string) (*TwitchConfigs, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// json data
	var config TwitchConfigs
	// unmarshall it
	err = json.Unmarshal(fdata, &config)
	if err != nil {
		return nil, err
	}

	config.TwitchIRL = strings.Replace(config.TwitchIRL, "irc://", "", 1)

	return &config, nil
}

func (config *TwitchConfigs) Reload(filepath string) bool {
	// parse config gile
	newConfig, err := Load(filepath)
	if err != nil {
		return false
	}
	changes := false
	if newConfig.TwitchIRL != config.TwitchIRL {
		config.TwitchIRL = newConfig.TwitchIRL
		changes = true
	}
	if newConfig.Channel != config.Channel {
		config.Channel = newConfig.Channel
		changes = true
	}
	if newConfig.Debug != config.Debug {
		config.Debug = newConfig.Debug
		changes = true
	}
	// check for changes in Commands
	for _, nwCfg := range newConfig.Commands {
		new := true
		for i, cfg := range config.Commands {
			if nwCfg.Name == cfg.Name {
				new = false
				if nwCfg.Enable != cfg.Enable {
					config.Commands[i].Enable = nwCfg.Enable
					changes = true
				}
				if nwCfg.Trigger != cfg.Trigger {
					config.Commands[i].Trigger = nwCfg.Trigger
					changes = true
				}
				break
			}
		}
		if new {
			config.Commands = append(config.Commands, nwCfg)
			changes = true
		}
	}
	// check for changes in Filters
	for _, nwCfg := range newConfig.Filters {
		new := true
		for i, cfg := range config.Filters {
			if nwCfg.Name == cfg.Name {
				new = false
				if nwCfg.Enable != cfg.Enable {
					config.Filters[i].Enable = nwCfg.Enable
					changes = true
				}
				if nwCfg.Pattern != cfg.Pattern {
					config.Filters[i].Pattern = nwCfg.Pattern
					changes = true
				}
				break
			}
		}
		if new {
			config.Filters = append(config.Filters, nwCfg)
			changes = true
		}
	}
	return changes
}
