package plugins

import (
	"log"
	"time"
)

type CommandConfig struct {
	Name    string
	Enable  bool
	Trigger string
}

type Command interface {
	Execute(string) string
}

type CommandMap map[string]Command

// Time Command
type timeComand struct {
}

func (cmd *timeComand) Execute(_ string) string {
	return time.Now().String()
}

// Uptime Command
type uptimeCommand struct {
	start time.Time
}

func (cmd *uptimeCommand) Execute(_ string) string {
	duration := time.Since(cmd.start)
	// format duration
	return time.Time{}.Add(duration).Format("15:04:05")
}

func LoadCommands(configs []CommandConfig) CommandMap {
	output := make(CommandMap, 0)

	for _, config := range configs {
		if config.Enable {
			switch config.Name {
			case "time":
				output[config.Trigger] = &timeComand{}
			case "uptime":
				output[config.Trigger] = &uptimeCommand{start: time.Now()}
			default:
				log.Printf("Unknowed command: %s\n", config.Name)
			}
		}
	}

	return output
}
