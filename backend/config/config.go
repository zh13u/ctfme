package config

import (
	"os"
	"strconv"
)

var TeamMode bool
var DynamicScoreEnabled bool
var DynamicScoreDecay int
var DynamicScoreMin int

// InitConfig initializes configuration after .env is loaded
func InitConfig() {
	TeamMode = os.Getenv("CTF_MODE") == "team"
	DynamicScoreEnabled = os.Getenv("DYNAMIC_SCORE_ENABLED") == "true"
	DynamicScoreDecay = getEnvAsInt("DYNAMIC_SCORE_DECAY", 10)
	DynamicScoreMin = getEnvAsInt("DYNAMIC_SCORE_MIN", 50)
}

// Helper function to get environment variable as int with default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Debug function to print current config
func PrintConfig() {
	println("=== CTF Configuration ===")
	println("CTF_MODE:", os.Getenv("CTF_MODE"))
	println("TeamMode:", TeamMode)
	println("DYNAMIC_SCORE_ENABLED:", os.Getenv("DYNAMIC_SCORE_ENABLED"))
	println("DynamicScoreEnabled:", DynamicScoreEnabled)
	println("DynamicScoreDecay:", DynamicScoreDecay)
	println("DynamicScoreMin:", DynamicScoreMin)
	println("========================")
}
