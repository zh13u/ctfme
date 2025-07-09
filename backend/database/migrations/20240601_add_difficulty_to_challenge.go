package migrations

import (
	"ctfme/database"
	"log"
)

func AddDifficultyToChallenge() {
	log.Println("Adding difficulty field to challenges table...")
	// Add column to challenges table
	err := database.DB.Exec("ALTER TABLE challenges ADD COLUMN IF NOT EXISTS difficulty VARCHAR(32) NOT NULL DEFAULT 'Easy';").Error
	if err != nil {
		log.Printf("Error adding difficulty column: %v", err)
	} else {
		log.Println("Difficulty field added to challenges table successfully")
	}
}
