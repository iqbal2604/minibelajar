package connection

import (
	"fmt"
	"log"
	"mini/config"
	"mini/database/ent"

	_ "github.com/lib/pq"
)

var db *ent.Client

func DB() *ent.Client {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.DB.Host, config.DB.Port, config.DB.User,
		config.DB.Pass, config.DB.Name, config.DB.SSL, config.DB.Timezone,
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return nil
	}

	db = client
	log.Println("Database connected")
	return client
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
