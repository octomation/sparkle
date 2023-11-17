package config

import "github.com/google/uuid"

type Server struct {
	Name    string `toml:"name" json:"name"`
	Service struct {
		License uuid.UUID `toml:"license" json:"license"`
	} `toml:"service" json:"service"`
	Sparkle struct {
		Path string `toml:"path" json:"path"`
		File string `toml:"file" json:"file"`
	} `toml:"sparkle" json:"sparkle"`
}
