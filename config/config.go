package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config func to get env value

func Config(key string, isFilePath bool) string {

	relPath := ".." + string(os.PathSeparator)
	err := godotenv.Load(filepath.Join(relPath, ".env"))
	if err != nil {
		fmt.Println("Error loading .env file", err)

	}

	//handle os specific path-separation
	value := os.Getenv(key)
	if osSpecificValue := filePathHandler(value); osSpecificValue != "" && isFilePath {
		value = filePathHandler(value)
	}

	return value
}

// handle separator by OS
func filePathHandler(path string) string {
	sep := string(os.PathSeparator)
	path = strings.ReplaceAll(path, "#", sep)
	return path
}
