package conf

import (
	"github.com/peerless6372/lizh/env"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func LoadConf(filename, subConf string, s interface{}) {
	var path string
	path = filepath.Join(env.GetConfDirPath(), subConf, filename)

	if yamlFile, err := os.ReadFile(path); err != nil {
		panic(filename + " get error: %v " + err.Error())
	} else if err = yaml.Unmarshal(yamlFile, s); err != nil {
		panic(filename + " unmarshal error: %v" + err.Error())
	}
}
