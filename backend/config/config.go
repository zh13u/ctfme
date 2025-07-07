package config

import (
	"ctfme/database"
	"ctfme/models"
	"log"
)

var TeamMode bool
var DynamicScoreEnabled bool
var DynamicScoreDecay int
var DynamicScoreMin int

// InitConfig initializes configuration from database
func InitConfig() {
	var config models.Setup
	if err := database.DB.First(&config).Error; err != nil {
		// If no config exists, create default config
		log.Println("No config found in database, creating default config...")
		defaultConfig := models.Setup{
			CTFMode:             "user",
			DynamicScoreEnabled: false,
			DynamicScoreDecay:   10,
			DynamicScoreMin:     50,
		}
		if err := database.DB.Create(&defaultConfig).Error; err != nil {
			log.Fatal("Failed to create default config: ", err)
		}
		config = defaultConfig
	}

	// Load config from database
	TeamMode = config.CTFMode == "team"
	DynamicScoreEnabled = config.DynamicScoreEnabled
	DynamicScoreDecay = config.DynamicScoreDecay
	DynamicScoreMin = config.DynamicScoreMin
}

// ReloadConfig reloads configuration from database
func ReloadConfig() error {
	var config models.Setup
	if err := database.DB.First(&config).Error; err != nil {
		return err
	}

	TeamMode = config.CTFMode == "team"
	DynamicScoreEnabled = config.DynamicScoreEnabled
	DynamicScoreDecay = config.DynamicScoreDecay
	DynamicScoreMin = config.DynamicScoreMin

	return nil
}

// Debug function to print current config
func PrintConfig() {
	println("=== CTF Configuration ===")
	println("CTF_MODE:", getCTFModeString())
	println("TeamMode:", TeamMode)
	println("DynamicScoreEnabled:", DynamicScoreEnabled)
	println("DynamicScoreDecay:", DynamicScoreDecay)
	println("DynamicScoreMin:", DynamicScoreMin)
	println("========================")
}

// Helper function to get CTF mode as string
func getCTFModeString() string {
	if TeamMode {
		return "team"
	}
	return "user"
}
