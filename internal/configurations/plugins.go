package configurations

type CommandConfig struct {
	Name    string `json:"name"`
	Enable  bool   `json:"enable"`
	Trigger string `json:"trigger"`
}

type FilterConfig struct {
	Name    string `json:"name"`
	Enable  bool   `json:"enable"`
	Pattern string `json:"pattern"`
}
