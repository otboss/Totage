CREATE SCHEMA IF NOT EXISTS _user;

CREATE TABLE IF NOT EXISTS _user.users
(
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT UNIQUE,
    password_hash VARCHAR(128),
    password_last_changed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    google_user_id TEXT,
    oauth_refresh_token_encrypted TEXT,
    last_ip_address TEXT NOT NULL,
    terms_and_condition_accepted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS _user.system_administrators
(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES _user.users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _user.verification_sessions
(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL,
    verification_codes_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES _user.users(id) ON DELETE CASCADE,
    CONSTRAINT fk_verification_code FOREIGN KEY (verification_codes_id) REFERENCES shared.verification_codes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _user.email_verification_codes
(
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    verification_codes_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_verification_code FOREIGN KEY (verification_codes_id) REFERENCES shared.verification_codes(id) ON DELETE CASCADE
);