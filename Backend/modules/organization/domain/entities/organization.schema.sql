CREATE SCHEMA IF NOT EXISTS organization;
CREATE TYPE organization.bank_account_ownership AS ENUM ('PERSONAL', 'BUSINESS');
CREATE TYPE organization.banK_account_type AS ENUM ('SAVINGS', 'CHEQUING');
CREATE TYPE organization.preferred_liquid_unit AS ENUM ('GALLON', 'LITRE');
CREATE TYPE organization.preferred_solid_unit AS ENUM ('EACH', 'DOZEN', 'KILOGRAM');

CREATE TABLE IF NOT EXISTS organization.organizations
(
    id BIGSERIAL PRIMARY KEY,
    avatar_url TEXT,
    instagram_handle TEXT,
    -- instagram_handle_verification_code VARCHAR(6) NOT NULL DEFAULT UPPER(SUBSTRING(gen_random_uuid()::TEXT FROM 1 FOR 6)),
    -- instagram_handle_verification_code_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    organization_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL,
    -- email_addresses TEXT[] NOT NULL CHECK (email_addresses IS NOT NULL AND array_length(email_addresses, 1) > 0 AND array_length(email_addresses, 1) <= 3),
    country shared.country_code NOT NULL,
    currency shared.currency NOT NULL,
    address TEXT NOT NULL,
    address_geog GEOGRAPHY(POINT, 4326) NOT NULL,
    duty_free_km FLOAT NOT NULL DEFAULT 0,
    max_order_queue_size SMALLINT,
    -- min_selection_time_in_milliseconds BIGINT NOT NULL DEFAULT 540000,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5) DEFAULT 5,
    -- is_water_service_supported BOOLEAN NOT NULL DEFAULT FALSE,
    -- is_rental_vehicle_service_supported BOOLEAN NOT NULL DEFAULT FALSE,
    -- is_freight_service_supported BOOLEAN NOT NULL DEFAULT FALSE,
    -- is_block_service_supported BOOLEAN NOT NULL DEFAULT FALSE,
    -- is_cooconut_service_supported BOOLEAN NOT NULL DEFAULT FALSE,
    -- is_egg_service_supported BOOLEAN NOT NULL DEFAULT FALSE,
    is_email_and_phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_submitted_for_approval BOOLEAN NOT NULL DEFAULT FALSE,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_organization_address_geog ON organization.organizations USING GIST(address_geog);

CREATE TABLE IF NOT EXISTS organization.instagram_handle_verification
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    instagram_handle_verification_code VARCHAR(6) NOT NULL DEFAULT UPPER(SUBSTRING(gen_random_uuid()::TEXT FROM 1 FOR 6)),
    instagram_handle_verification_code_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    failed_attempts SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.organization_approvals (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL,
    system_administrator_id BIGSERIAL NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_sys_admin FOREIGN KEY (system_administrator_id) REFERENCES _user.system_administrators(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.organization_mail_list
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL,
    email TEXT NOT NULL,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE,
    CONSTRAINT unique_organization_email_pair UNIQUE (organization_id, email)
);

CREATE TABLE IF NOT EXISTS organization.organization_verification_code
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    phone_verification_codes_id UUID NOT NULL UNIQUE,
    email_verification_codes_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_phone_verification_code FOREIGN KEY (phone_verification_codes_id) REFERENCES shared.verification_codes(id) ON DELETE CASCADE,
    CONSTRAINT fk_email_verification_code FOREIGN KEY (email_verification_codes_id) REFERENCES shared.verification_codes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.accounts 
(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGSERIAL NOT NULL,
    organization_id BIGSERIAL NOT NULL,
    is_invitation_accepted BOOLEAN NOT NULL DEFAULT FALSE,
    is_admin BOOLEAN NOT NULL,
    is_customer_support BOOLEAN NOT NULL,
    is_driver BOOLEAN NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES _user.users(id) ON DELETE CASCADE,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE,
    UNIQUE (user_id, organization_id)
);

-- CREATE TABLE IF NOT EXISTS organization.account_verification_tokens
-- (
--     id BIGSERIAL PRIMARY KEY,
--     account_id BIGSERIAL NOT NULL UNIQUE,
--     verification_token VARCHAR(255) NOT NULL UNIQUE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS organization.timeslots
-- (
--     id BIGSERIAL PRIMARY KEY,
--     organization_id BIGSERIAL NOT NULL,
--     start_hh   SMALLINT NOT NULL CHECK (start_hh BETWEEN 0 AND 23),
--     start_mm SMALLINT NOT NULL CHECK (start_mm = 0 OR start_mm = 30),
--     end_hh   SMALLINT NOT NULL CHECK (end_hh BETWEEN 0 AND 23),
--     end_mm SMALLINT NOT NULL CHECK (end_mm = 0 OR end_mm = 30),
--     wkdays shared.weekday[] NOT NULL CHECK ( array_length(wkdays, 1) > 0 ),
--     average_job_duration BIGINT NOT NULL,
--     order_type shared.order_type NOT NULL,
--     UNIQUE(organization_id, order_type),
--     CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
-- );

CREATE TABLE IF NOT EXISTS organization.potable_water_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_liquid_unit NOT NULL DEFAULT 'GALLON',
    water_price_min FLOAT NOT NULL,
    water_price_per_km FLOAT NOT NULL,
    water_price FLOAT NOT NULL,
    bottle_cost FLOAT,
    available_supply FLOAT NOT NULL,
    min_order_quantity FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.irrigation_water_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_liquid_unit NOT NULL DEFAULT 'GALLON',
    water_price_min FLOAT NOT NULL,
    water_price_per_km FLOAT NOT NULL,
    water_price FLOAT NOT NULL,
    barrel_cost FLOAT,
    available_supply FLOAT NOT NULL,
    min_order_quantity FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.block_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_solid_unit NOT NULL DEFAULT 'EACH',
    block_price_min FLOAT NOT NULL,
    block_price FLOAT NOT NULL,
    block_price_per_km FLOAT NOT NULL,
    available_supply FLOAT NOT NULL,
    min_order_quantity FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.ice_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_solid_unit NOT NULL DEFAULT 'KILOGRAM',
    ice_price_min FLOAT NOT NULL,
    ice_price FLOAT NOT NULL,
    ice_price_per_km FLOAT NOT NULL,
    available_supply FLOAT NOT NULL,
    min_order_quantity FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.coconut_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_solid_unit NOT NULL DEFAULT 'EACH',
    coconut_price_min FLOAT NOT NULL,
    coconut_price FLOAT NOT NULL,
    coconut_price_per_km  FLOAT NOT NULL,
    available_supply FLOAT NOT NULL,
    min_order_quantity FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.egg_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_solid_unit NOT NULL DEFAULT 'EACH',
    egg_price_min FLOAT NOT NULL,
    egg_price FLOAT NOT NULL,
    egg_price_per_km FLOAT NOT NULL,
    available_supply FLOAT NOT NULL,
    min_order_quantity FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.freight_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    freight_price_min FLOAT NOT NULL,
    solid_waste_dump_site_geog GEOGRAPHY(POINT, 4326),
    liquid_waste_dump_site_geog GEOGRAPHY(POINT, 4326),
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.charcoal_sales_config
(
    -- TODO:
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_solid_unit NOT NULL DEFAULT 'EACH',
    charcoal_price_min FLOAT NOT NULL,
    charcoal_price FLOAT NOT NULL,
    charcoal_price_per_km  FLOAT NOT NULL,
    available_supply FLOAT NOT NULL,
    min_order_quantity FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.honey_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_liquid_unit NOT NULL DEFAULT 'LITRE',
    honey_price_min FLOAT NOT NULL,
    honey_price FLOAT NOT NULL,
    honey_price_per_km  FLOAT NOT NULL,
    available_supply FLOAT NOT NULL,
    min_order_quantity FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS organization.sales_units
-- (
--     id BIGSERIAL PRIMARY KEY,
--     organization_id BIGSERIAL NOT NULL,
--     order_type shared.order_type NOT NULL,
--     quantity SMALLINT NOT NULL,
--     CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
-- );

CREATE TABLE IF NOT EXISTS organization.paypal_credentials
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    email TEXT NOT NULL,
    subscription_id TEXT NOT NULL,
    rest_sdk_client_id_encrypted TEXT,
    rest_sdk_secret_encrypted TEXT,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE,
    CHECK (
        (rest_sdk_client_id_encrypted IS NULL AND rest_sdk_secret_encrypted IS NULL) OR
        (rest_sdk_client_id_encrypted IS NOT NULL AND rest_sdk_secret_encrypted IS NOT NULL)
    )
);

CREATE TABLE IF NOT EXISTS organization.bank_accounts 
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL,
    bank TEXT NOT NULL,
    bank_branch TEXT NOT NULL,
    account_name TEXT NOT NULL,
    type organization.banK_account_type NOT NULL,
    ownership organization.bank_account_ownership NOT NULL,
    account_number TEXT NOT NULL,
    swift_code TEXT NOT NULL,
    routing_number TEXT,
    recipient_address TEXT,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);
