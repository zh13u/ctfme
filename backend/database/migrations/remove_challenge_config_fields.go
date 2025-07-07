package migrations

import (
	"ctfme/database"
	"log"
)

func RemoveChallengeConfigFields() {
	log.Println("Removing config fields from challenges table...")

	// Drop columns from challenges table
	err := database.DB.Exec("ALTER TABLE challenges DROP COLUMN IF EXISTS dynamic_score").Error
	if err != nil {
		log.Printf("Error dropping dynamic_score column: %v", err)
	}

	err = database.DB.Exec("ALTER TABLE challenges DROP COLUMN IF EXISTS min_score").Error
	if err != nil {
		log.Printf("Error dropping min_score column: %v", err)
	}

	err = database.DB.Exec("ALTER TABLE challenges DROP COLUMN IF EXISTS decay").Error
	if err != nil {
		log.Printf("Error dropping decay column: %v", err)
	}

	log.Println("Config fields removed from challenges table successfully")
}
