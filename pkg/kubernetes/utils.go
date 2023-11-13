package kubernetes

import (
	"log"
	"os"
	"path"
)

func getResourcesPath(name string) string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	rPath := path.Join(currentPath, "../../deployments", name)

	return rPath
}

func getYAMLfiles(name string) [][]byte {
	var yamlFiles [][]byte
	rPath := getResourcesPath(name)

	files, err := os.ReadDir(rPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := os.ReadFile(path.Join(rPath, "/", file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, data)
	}

	return yamlFiles
}
