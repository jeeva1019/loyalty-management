package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// TomlMap stores the decoded TOML configurations.
//   - Key: filename without extension, Value: map of key-values in that TOML file.
var TomlMap = make(map[string]map[string]string)

// TomlInit reads all .toml files from the "./settings" directory and decodes them into TomlMap.
func TomlInit() {
	folderPath := "./settings"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("❌ Failed to read settings folder: %v", err)
	}

	for _, file := range files {
		// Skip directories and non-TOML files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
			continue
		}

		filePath := filepath.Join(folderPath, file.Name())
		innerMap := make(map[string]string)

		// Decode TOML file into map[string]string
		if _, err := toml.DecodeFile(filePath, &innerMap); err != nil {
			log.Fatalf("❌ Error loading TOML file %s: %v", file.Name(), err)
		}

		// Use filename without extension as key
		key := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		TomlMap[key] = innerMap
	}

	fmt.Println("✅ TOML configuration loaded successfully")
}

// GetTomlStr returns the string value for the given filename and key.
//   - Returns empty string if file or key is not found.
func GetTomlStr(fileName, key string) string {
	if valMap, ok := TomlMap[fileName]; ok {
		if val, exists := valMap[key]; exists {
			return val
		}
		log.Printf("⚠️ TOML key not found: [%s][%s]", fileName, key)
	} else {
		log.Printf("⚠️ TOML file not found: %s %s", fileName, key)
	}
	return ""
}

// GetTomlMap returns the entire key-value map for the given TOML file.
//   - Returns nil if file is not found.
func GetTomlMap(fileName string) map[string]string {
	if valMap, ok := TomlMap[fileName]; ok {
		return valMap
	}
	log.Printf("⚠️ TOML file not found: %s", fileName)
	return nil
}
