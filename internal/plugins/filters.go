package plugins

import (
	"log"
	"strings"
)

type FilterConfig struct {
	Name    string `json:"name"`
	Enable  bool   `json:"enable"`
	Pattern string `json:"pattern"`
}

type Filter interface {
	Apply(string) bool
}

type FilterMap map[string]Filter

// Starts With Filter
type startsFilter struct {
	Pattern string
}

func (filter *startsFilter) Apply(data string) bool {
	return strings.HasPrefix(data, filter.Pattern)
}

// Contains Filter
type ContainsFilter struct {
	Pattern string
}

func (filter *ContainsFilter) Apply(data string) bool {
	return strings.Contains(data, filter.Pattern)
}

func LoadFilter(filterConfigs []FilterConfig) FilterMap {
	output := make(FilterMap)
	for _, config := range filterConfigs {
		if config.Enable {
			switch config.Name {
			case "StartsWith":
				output[config.Name] = &startsFilter{Pattern: config.Pattern}
			case "Contains":
				output[config.Name] = &ContainsFilter{Pattern: config.Pattern}
			default:
				log.Printf("Unknowed filter: %s\n", config.Name)
			}
		}
	}
	return output
}
