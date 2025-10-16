package database

import (
	"database/sql"
	"fmt"
)

// GetSetting retrieves a setting value by key
func GetSetting(key string) (string, error) {
	var value string
	query := `SELECT value FROM settings WHERE key = ?`
	err := db.QueryRow(query, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("setting not found: %s", key)
		}
		return "", fmt.Errorf("database error: %w", err)
	}
	return value, nil
}

// SetSetting sets or updates a setting
func SetSetting(key, value string) error {
	query := `INSERT INTO settings (key, value) VALUES (?, ?)
			  ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP`
	_, err := db.Exec(query, key, value, value)
	if err != nil {
		return fmt.Errorf("failed to set setting: %w", err)
	}
	return nil
}

// DeleteSetting deletes a setting by key
func DeleteSetting(key string) error {
	query := `DELETE FROM settings WHERE key = ?`
	result, err := db.Exec(query, key)
	if err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("setting not found: %s", key)
	}

	return nil
}

// GetAllSettings retrieves all settings
func GetAllSettings() (map[string]string, error) {
	settings := make(map[string]string)
	query := `SELECT key, value FROM settings`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve settings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings[key] = value
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating settings: %w", err)
	}

	return settings, nil
}
