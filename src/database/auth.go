package database

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AdminCredentials represents admin user credentials
type AdminCredentials struct {
	ID           int
	Username     string
	PasswordHash string
	Token        string
	CreatedAt    time.Time
	LastLogin    *time.Time
}

// CreateAdmin creates a new admin user
func CreateAdmin(username, password, token string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `INSERT INTO admins (username, password_hash, token) VALUES (?, ?, ?)`
	_, err = db.Exec(query, username, hashedPassword, token)
	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	return nil
}

// ValidateAdminCredentials validates username and password
func ValidateAdminCredentials(username, password string) (*AdminCredentials, error) {
	var admin AdminCredentials
	var lastLogin sql.NullTime

	query := `SELECT id, username, password_hash, token, created_at, last_login FROM admins WHERE username = ?`
	err := db.QueryRow(query, username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.Token,
		&admin.CreatedAt,
		&lastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid username or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if lastLogin.Valid {
		admin.LastLogin = &lastLogin.Time
	}

	// Update last login
	_, _ = db.Exec(`UPDATE admins SET last_login = CURRENT_TIMESTAMP WHERE id = ?`, admin.ID)

	return &admin, nil
}

// ValidateAdminToken validates an admin token
func ValidateAdminToken(token string) (*AdminCredentials, error) {
	var admin AdminCredentials
	var lastLogin sql.NullTime

	query := `SELECT id, username, password_hash, token, created_at, last_login FROM admins WHERE token = ?`
	err := db.QueryRow(query, token).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.Token,
		&admin.CreatedAt,
		&lastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid token")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	if lastLogin.Valid {
		admin.LastLogin = &lastLogin.Time
	}

	return &admin, nil
}

// AdminExists checks if any admin user exists
func AdminExists() (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM admins`).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateAdminPassword updates an admin's password
func UpdateAdminPassword(username, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `UPDATE admins SET password_hash = ? WHERE username = ?`
	result, err := db.Exec(query, hashedPassword, username)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("admin not found")
	}

	return nil
}

// UpdateAdminToken updates an admin's token
func UpdateAdminToken(username, newToken string) error {
	query := `UPDATE admins SET token = ? WHERE username = ?`
	result, err := db.Exec(query, newToken, username)
	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("admin not found")
	}

	return nil
}
