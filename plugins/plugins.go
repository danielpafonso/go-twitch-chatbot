package plugins

import (
	"errors"
	"fmt"
	"path"
	"plugin"

	"twitch-chatbot/internal/configurations"
)

// Functions that command plugins must implement
type Command interface {
	Initiate(args ...interface{}) error
	Execute(args ...interface{}) (string, error)
}

// Functions that filter plugins must implement
type Filter interface {
	Initiate(args ...interface{}) error
	Apply(line string) (string, error)
}

type CommandMap map[string]Command
type FilterMap map[string]Filter

// LoadPluginsCommands load and configure commands plugins
func LoadPluginsCommands(configs []configurations.CommandConfig, pluginFolder string) (CommandMap, error) {
	commandMap := make(CommandMap, 0)
	for _, config := range configs {
		if !config.Enable {
			continue
		}
		fmt.Printf("loading %s with trigger %s\n", config.Name, config.Trigger)
		pluginObject := path.Join("", pluginFolder, config.Name+".so")
		// read file
		plug, err := plugin.Open(pluginObject)
		if err != nil {
			return nil, err
		}
		// look up symbol
		symPlug, err := plug.Lookup("Command")
		if err != nil {
			return nil, err
		}
		// asset loading correct
		cmd, ok := symPlug.(Command)
		if !ok {
			return nil, errors.New("unexpected type from module symbol")
		}
		// initiate plugin
		err = cmd.Initiate()
		if err != nil {
			return nil, err
		}
		// add command to map
		commandMap[config.Trigger] = cmd
	}

	return commandMap, nil
}

// LoadPluginsFilter load and configure commands plugins
func LoadPluginsFilter(configs configurations.FilterConfig, pluginFolder string) (FilterMap, error) {
	filterMap := make(FilterMap, 0)

	return filterMap, nil
}
