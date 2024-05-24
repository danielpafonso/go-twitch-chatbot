package configurations

import (
	"encoding/json"
	"os"
	"strings"
)

type TwitchConfigs struct {
	TwitchIRL string `json:"twicthIrc"`
	Channel   string `json:"channel"`
	Debug     bool   `json:"debug"`
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
	return changes
}
