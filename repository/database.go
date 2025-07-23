package repository

import (
	"database/sql"
	"fmt"
)

func InitDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		customer_id VARCHAR(255) NOT NULL,
		product_id VARCHAR(255) NOT NULL,
		quantity INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	return nil
}
