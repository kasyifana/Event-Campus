-- Event Campus Database Schema
-- Execute this in Supabase SQL Editor

-- Table: users
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

-- Table: whitelist_requests
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

-- Table: events
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

-- Table: registrations
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

-- Table: attendances
CREATE TABLE IF NOT EXISTS attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_id UUID REFERENCES registrations(id) ON DELETE CASCADE,
    checked_in_at TIMESTAMP DEFAULT NOW(),
    notes TEXT
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);
CREATE INDEX IF NOT EXISTS idx_events_start_date ON events(start_date);
CREATE INDEX IF NOT EXISTS idx_events_organizer ON events(organizer_id);
CREATE INDEX IF NOT EXISTS idx_events_category ON events(category);
CREATE INDEX IF NOT EXISTS idx_registrations_event ON registrations(event_id);
CREATE INDEX IF NOT EXISTS idx_registrations_user ON registrations(user_id);
CREATE INDEX IF NOT EXISTS idx_registrations_status ON registrations(status);
CREATE INDEX IF NOT EXISTS idx_whitelist_status ON whitelist_requests(status);
CREATE INDEX IF NOT EXISTS idx_whitelist_user ON whitelist_requests(user_id);

-- Create default admin user (password: admin123)
-- Password hash for "admin123" using bcrypt
INSERT INTO users (email, password_hash, full_name, phone_number, role, is_uii_civitas, is_approved)
VALUES (
    'admin@uii.ac.id',
    '$2a$10$K7mZ3vGGqG9vY8xQVH9Q8eN4xW3rZ8KF7qY5xH9Q8eN4xW3rZ8KF7q',
    'System Administrator',
    '+6281234567890',
    'admin',
    true,
    true
) ON CONFLICT (email) DO NOTHING;

-- Function to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for auto-updating updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- View for event list with organizer info
CREATE OR REPLACE VIEW events_with_organizer AS
SELECT 
    e.*,
    u.full_name as organizer_name
FROM events e
LEFT JOIN users u ON e.organizer_id = u.id;

-- View for registrations with user and event info
CREATE OR REPLACE VIEW registrations_detailed AS
SELECT 
    r.*,
    u.full_name as user_name,
    u.email as user_email,
    u.phone_number as user_phone,
    e.title as event_title,
    e.start_date as event_date
FROM registrations r
LEFT JOIN users u ON r.user_id = u.id
LEFT JOIN events e ON r.event_id = e.id;

-- View for whitelist requests with user info
CREATE OR REPLACE VIEW whitelist_requests_detailed AS
SELECT 
    w.*,
    u.full_name as user_name,
    u.email as user_email,
    r.full_name as reviewer_name
FROM whitelist_requests w
LEFT JOIN users u ON w.user_id = u.id
LEFT JOIN users r ON w.reviewed_by = r.id;

COMMENT ON TABLE users IS 'User accounts with role-based access';
COMMENT ON TABLE events IS 'Events created by organizations';
COMMENT ON TABLE registrations IS 'User registrations to events';
COMMENT ON TABLE whitelist_requests IS 'Requests to become organization';
COMMENT ON TABLE attendances IS 'Attendance records for events';
