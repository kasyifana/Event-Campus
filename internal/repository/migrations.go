package repository

import (
	"context"
	"database/sql"
	"log"
)

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB) error {
	ctx := context.Background()

	// Create users table
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(255) NOT NULL,
			phone_number VARCHAR(20) NOT NULL,
			role VARCHAR(20) NOT NULL CHECK (role IN ('mahasiswa', 'organisasi', 'admin')),
			is_uii_civitas BOOLEAN DEFAULT FALSE,
			is_approved BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}
	log.Println("✅ Table 'users' ready")

	// Create whitelist_requests table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS whitelist_requests (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			organization_name VARCHAR(255) NOT NULL,
			document_path VARCHAR(500) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
			admin_notes TEXT,
			submitted_at TIMESTAMP DEFAULT NOW(),
			reviewed_at TIMESTAMP,
			reviewed_by UUID REFERENCES users(id)
		);
	`)
	if err != nil {
		return err
	}
	log.Println("✅ Table 'whitelist_requests' ready")

	// Create events table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			organizer_id UUID REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			category VARCHAR(50) NOT NULL CHECK (category IN ('seminar', 'workshop', 'lomba', 'konser')),
			event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('online', 'offline')),
			location VARCHAR(255),
			zoom_link VARCHAR(500),
			poster_path VARCHAR(500),
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP NOT NULL,
			registration_deadline TIMESTAMP NOT NULL,
			max_participants INT NOT NULL,
			current_participants INT DEFAULT 0,
			is_uii_only BOOLEAN DEFAULT FALSE,
			status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'ongoing', 'completed', 'cancelled')),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}
	log.Println("✅ Table 'events' ready")

	// Create registrations table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS registrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID REFERENCES events(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			status VARCHAR(20) DEFAULT 'registered' CHECK (status IN ('registered', 'waitlist', 'cancelled', 'attended')),
			registered_at TIMESTAMP DEFAULT NOW(),
			cancelled_at TIMESTAMP,
			reminder_sent BOOLEAN DEFAULT FALSE,
			UNIQUE(event_id, user_id)
		);
	`)
	if err != nil {
		return err
	}
	log.Println("✅ Table 'registrations' ready")

	// Create attendances table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS attendances (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			registration_id UUID REFERENCES registrations(id) ON DELETE CASCADE,
			checked_in_at TIMESTAMP DEFAULT NOW(),
			notes TEXT
		);
	`)
	if err != nil {
		return err
	}
	log.Println("✅ Table 'attendances' ready")

	// Create indexes
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`)
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`)
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);`)
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_events_start_date ON events(start_date);`)
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_events_organizer ON events(organizer_id);`)
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_events_category ON events(category);`)
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_registrations_event ON registrations(event_id);`)
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_registrations_user ON registrations(user_id);`)
	db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_registrations_status ON registrations(status);`)
	log.Println("✅ Indexes created")

	// Insert default admin if not exists
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", "admin@uii.ac.id").Scan(&count)
	if err == nil && count == 0 {
		// Password: admin123 (bcrypt hash)
		_, err = db.ExecContext(ctx, `
			INSERT INTO users (email, password_hash, full_name, phone_number, role, is_uii_civitas, is_approved)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, "admin@uii.ac.id",
			"$2a$10$K7mZ3vGGqG9vY8xQVH9Q8eN4xW3rZ8KF7qY5xH9Q8eN4xW3rZ8KF7q",
			"System Administrator",
			"+6281234567890",
			"admin",
			true,
			true,
		)
		if err == nil {
			log.Println("✅ Default admin created (admin@uii.ac.id / admin123)")
		}
	}

	log.Println("🎉 Database schema initialized successfully!")
	return nil
}
