package server

import (
	"path"
)

type ServerConfigurable struct {
	Network struct {
		Port uint16
		IP   string
	}

	Computing struct {
		MaxWorkers uint
	}
}

func ReadConfigFromFile(configPath string) (ServerConfigurable, error) {
	workPath := path.Clean(configPath)
	extension := path.Ext(workPath)

	switch extension {
	case ".cfg", ".config", ".toml":
		break
	default:
		return ServerConfigurable{}, nil
	}

	return ServerConfigurable{}, nil
}
