package migrations

import (
	"ctfme/database"
	"log"
)

func AddCurrentPointsToChallenge() {
	log.Println("Adding current_points field to challenges table...")
	err := database.DB.Exec("ALTER TABLE challenges ADD COLUMN IF NOT EXISTS current_points INTEGER NOT NULL DEFAULT 0;").Error
	if err != nil {
		log.Printf("Error adding current_points column: %v", err)
	} else {
		log.Println("current_points field added to challenges table successfully")
	}
} 