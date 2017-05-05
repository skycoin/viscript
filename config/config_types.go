package config

type App struct {
	Path string   `yaml:"path"`
	Args []string `yaml:"args"`
	Desc string   `yaml:"desc"`
}

type Settings struct {
	VerboseInput bool `yaml:"verboseInput"`
}

type Config struct {
	Apps     map[string]App `yaml:"apps"`
	Settings Settings       `yaml:"settings"`
}
