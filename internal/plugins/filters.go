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
	GetPattern() string
}

type FilterMap map[string]Filter

// Starts With Filter
type startsFilter struct {
	Pattern string
}

func (filter *startsFilter) Apply(data string) bool {
	return strings.HasPrefix(data, filter.Pattern)
}
func (filter *startsFilter) GetPattern() string {
	return filter.Pattern
}

// Contains Filter
type containsFilter struct {
	Pattern string
}

func (filter *containsFilter) Apply(data string) bool {
	return strings.Contains(data, filter.Pattern)
}
func (filter *containsFilter) GetPattern() string {
	return filter.Pattern
}

func LoadFilter(filterConfigs []FilterConfig) FilterMap {
	output := make(FilterMap)
	for _, config := range filterConfigs {
		if config.Enable {
			switch config.Name {
			case "StartsWith":
				output[config.Name] = &startsFilter{Pattern: config.Pattern}
			case "Contains":
				output[config.Name] = &containsFilter{Pattern: config.Pattern}
			default:
				log.Printf("Unknowed filter: %s\n", config.Name)
			}
		}
	}
	return output
}
