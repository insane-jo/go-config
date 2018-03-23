package main

import (
	"os"
	"fmt"
	"path/filepath"
	"log"
	"encoding/json"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"encoding/xml"
)

var allowedExtensions = []string{"json","yaml","xml"}

func getCurrentCwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

type configFileType struct{
	FilePath string
	Extension string
}

func getConfigFiles() *[]configFileType {
	goenv := os.Getenv("GO_ENV")
	goinstance := os.Getenv("GO_APP_INSTANCE")
	configPrefix := os.Getenv("GO_CONFIG_PREFIX")
	configDirecroryName := os.Getenv("GO_CONFIG_DIR")

	if configPrefix == "" {
		configPrefix = "local-"
	}

	if configDirecroryName == "" {
		configDirecroryName = getCurrentCwd() + "/config"
	}

	filenames := []string{"default", configPrefix + "default"}

	if goenv != "" {
		filenames = append(filenames, goenv, configPrefix + goenv)
	}

	if goinstance != "" {
		if goenv != "" {
			filenames = append(filenames, goenv + "-" + goinstance, configPrefix + goenv + "-" + goinstance)
		} else {
			filenames = append(filenames, goinstance, configPrefix + goinstance)
		}
	}

	configFilenames := []configFileType{}
	for i := 0; i < len(filenames); i++ {
		for extensionId := 0; extensionId < len(allowedExtensions); extensionId++ {
			extension := allowedExtensions[extensionId]
			matches, err := filepath.Glob(configDirecroryName + "/" + filenames[i] + "." + extension)
			if err != nil {
				log.Fatal(err)
			}
			if len(matches) > 0 {
				for filenameIdx := 0; filenameIdx < len(matches); filenameIdx++ {
					configFilenames = append(configFilenames, configFileType{
						Extension:extension,
						FilePath:matches[filenameIdx],
					})
				}
				break
			}
		}
	}

	return &configFilenames
}

func ParseConfig(value interface{}) {
	configFiles := *getConfigFiles()

	for i := 0; i < len(configFiles); i++ {
		currentConfigFile := (configFiles)[i]
		filedata, err := ioutil.ReadFile(currentConfigFile.FilePath)
		if err != nil{
			log.Fatal(err)
		}
		if currentConfigFile.Extension == "json" {
			json.Unmarshal(filedata, value)
		} else if currentConfigFile.Extension == "yaml" {
			yaml.Unmarshal(filedata, value)
		} else if currentConfigFile.Extension == "xml" {
			xml.Unmarshal(filedata, value)
		}
	}
}

func main() {
	data := struct{
		Name string
		Host string
		Port int64
		LocalPrefix bool
	}{}
	ParseConfig(&data)

	fmt.Println(data)
}