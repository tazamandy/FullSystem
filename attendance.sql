CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    course VARCHAR(100),
    year_level VARCHAR(50),
    section VARCHAR(50),
    department VARCHAR(100),
    college VARCHAR(100),
    contact_number VARCHAR(20),
    address TEXT,
    qr_code_data TEXT,
    qr_type VARCHAR(50) DEFAULT 'student_id',
    qr_generated_at TIMESTAMP,
    role VARCHAR(50) NOT NULL DEFAULT 'student',
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMP,
    active_event_id INTEGER,
    original_qr_code_data TEXT,
    original_qr_type VARCHAR(50)
);


CREATE TABLE password_resets (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used BOOLEAN DEFAULT FALSE
);

-- Create indexes
CREATE INDEX idx_password_resets_email ON password_resets(email);
CREATE INDEX idx_password_resets_code ON password_resets(code);

CREATE TABLE pending_users (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    course VARCHAR(100),
    year_level VARCHAR(10),
    section VARCHAR(50),
    department VARCHAR(100),
    college VARCHAR(100),
    contact_number VARCHAR(20),
    address TEXT,
    verification_code VARCHAR(6) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP
);




ALTER TABLE users
ADD CONSTRAINT uni_users_student_id UNIQUE (student_id);

ALTER TABLE users
ADD CONSTRAINT uni_users_email UNIQUE (email);

-- Add profile_picture column to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS profile_picture TEXT;

-- Add profile_picture column to pending_users table
ALTER TABLE pending_users
ADD COLUMN IF NOT EXISTS profile_picture TEXT;

-- Create events table
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date DATE NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    location VARCHAR(255),
    course VARCHAR(100),
    section VARCHAR(50),
    year_level VARCHAR(50),
    department VARCHAR(100),
    college VARCHAR(100),
    created_by VARCHAR(255) NOT NULL,
    created_by_role VARCHAR(50) DEFAULT 'faculty',
    status VARCHAR(50) DEFAULT 'scheduled',
    is_active BOOLEAN DEFAULT TRUE,
    qr_code_data TEXT,
    tagged_courses TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for events table
CREATE INDEX IF NOT EXISTS idx_events_event_date ON events(event_date);
CREATE INDEX IF NOT EXISTS idx_events_created_by ON events(created_by);
CREATE INDEX IF NOT EXISTS idx_events_course ON events(course);
CREATE INDEX IF NOT EXISTS idx_events_section ON events(section);
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);
CREATE INDEX IF NOT EXISTS idx_events_is_active ON events(is_active);

-- Create attendances table
CREATE TABLE IF NOT EXISTS attendances (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    student_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'present',
    marked_at TIMESTAMP NOT NULL,
    marked_by VARCHAR(255),
    marked_by_role VARCHAR(50),
    method VARCHAR(50) DEFAULT 'qr_scan',
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    notes TEXT,

    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,
    check_in_status VARCHAR(50),
    check_out_status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(student_id) ON DELETE CASCADE
);

-- Create indexes for attendances table
CREATE INDEX IF NOT EXISTS idx_attendances_event_id ON attendances(event_id);
CREATE INDEX IF NOT EXISTS idx_attendances_student_id ON attendances(student_id);
CREATE INDEX IF NOT EXISTS idx_attendances_status ON attendances(status);
CREATE INDEX IF NOT EXISTS idx_attendances_marked_at ON attendances(marked_at);
CREATE INDEX IF NOT EXISTS idx_attendances_event_student ON attendances(event_id, student_id);


CREATE UNIQUE INDEX IF NOT EXISTS idx_attendances_unique_event_student 
ON attendances(event_id, student_id);


CREATE INDEX IF NOT EXISTS idx_attendances_check_in_time ON attendances(check_in_time);
CREATE INDEX IF NOT EXISTS idx_attendances_check_out_time ON attendances(check_out_time);


ALTER TABLE users
ADD COLUMN IF NOT EXISTS active_event_id INTEGER,
ADD COLUMN IF NOT EXISTS original_qr_code_data TEXT,
ADD COLUMN IF NOT EXISTS original_qr_type VARCHAR(50);


ALTER TABLE users
ADD CONSTRAINT fk_users_active_event
FOREIGN KEY (active_event_id) REFERENCES events(id) ON DELETE SET NULL;


CREATE INDEX IF NOT EXISTS idx_users_active_event_id ON users(active_event_id);

-- Create audit_logs table for tracking admin actions
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    action VARCHAR(255) NOT NULL,
    actor_id VARCHAR(255) NOT NULL,
    target_id VARCHAR(255),
    details TEXT,
    ip_address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for audit_logs table
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_id ON audit_logs(actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_target_id ON audit_logs(target_id);