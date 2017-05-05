package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var Global Config

const maxBufferSize = 4096

func Load(configFileName string) error {
	println("Loading configuration file:", configFileName)

	var err error

	file, err := os.Open(configFileName)
	if err != nil {
		return err
	}

	buffer := make([]byte, maxBufferSize)
	n, err := file.Read(buffer)
	if err != nil {
		return err
	}

	Global = Config{}

	err = yaml.Unmarshal(buffer[:n], &Global)
	if err != nil {
		return err
	}

	// This can be quickly used to see the contents of read config

	// fmt.Printf("[ Config ]\n")

	// for key, app := range Global.Apps {
	// 	fmt.Printf("[ App \"%s\" ]\n", key)
	// 	fmt.Printf("\tPath: %s\n", app.Path)
	// 	fmt.Printf("\tArgs: %v\n", app.Args)
	// 	fmt.Printf("\tDescription: %s\n\n", app.Desc)
	// }

	// fmt.Printf("Settings: %+v\n\n", Global.Settings)

	return nil
}

func DebugPrintInputEvents() bool {
	return Global.Settings.VerboseInput
}
