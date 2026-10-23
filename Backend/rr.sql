CREATE SCHEMA IF NOT EXISTS shared;

CREATE TYPE shared.currency AS ENUM (
    'AED', 'AFN', 'ALL', 'AMD', 'ANG', 'AOA', 'ARS', 'AUD', 'AWG', 'AZN', 'BAM', 'BBD',
    'BDT', 'BGN', 'BHD', 'BIF', 'BMD', 'BND', 'BOB', 'BOV', 'BRL', 'BSD', 'BTN', 'BWP',
    'BYN', 'BZD', 'CAD', 'CDF', 'CHE', 'CHF', 'CHW', 'CLF', 'CLP', 'CNY', 'COP', 'COU',
    'CRC', 'CUC', 'CUP', 'CVE', 'CZK', 'DJF', 'DKK', 'DOP', 'DZD', 'EGP', 'ERN', 'ETB',
    'EUR', 'FJD', 'FKP', 'GBP', 'GEL', 'GHS', 'GIP', 'GMD', 'GNF', 'GTQ', 'GYD', 'HKD',
    'HNL', 'HRK', 'HTG', 'HUF', 'IDR', 'ILS', 'IMP', 'INR', 'IQD', 'IRR', 'ISK', 'JEP',
    'JMD', 'JOD', 'JPY', 'KES', 'KGS', 'KHR', 'KMF', 'KPW', 'KRW', 'KWD', 'KYD', 'KZT',
    'LAK', 'LBP', 'LKR', 'LRD', 'LSL', 'LYD', 'MAD', 'MDL', 'MGA', 'MKD', 'MMK', 'MNT',
    'MOP', 'MRU', 'MUR', 'MVR', 'MWK', 'MXN', 'MXV', 'MYR', 'MZN', 'NAD', 'NGN', 'NIO',
    'NOK', 'NPR', 'NZD', 'OMR', 'PAB', 'PEN', 'PGK', 'PHP', 'PKR', 'PLN', 'PYG', 'QAR',
    'RON', 'RSD', 'RUB', 'RWF', 'SAR', 'SBD', 'SCR', 'SDG', 'SEK', 'SGD', 'SHP', 'SLL',
    'SOS', 'SRD', 'SSP', 'STN', 'SVC', 'SYP', 'SZL', 'THB', 'TJS', 'TMT', 'TND', 'TOP',
    'TRY', 'TTD', 'TWD', 'TZS', 'UAH', 'UGX', 'USD', 'USN', 'UYI', 'UYU', 'UYW', 'UZS',
    'VED', 'VES', 'VND', 'VUV', 'WST', 'XAF', 'XAG', 'XAU', 'XCD', 'XDR', 'XOF', 'XPD',
    'XPF', 'XPT', 'XSU', 'XTS', 'XUA', 'XXX', 'YER', 'ZAR', 'ZMW', 'ZWL'
);

CREATE TYPE shared.country_code AS ENUM (
    'AFG', 'ALA', 'ALB', 'DZA', 'ASM', 'AND', 'AGO', 'AIA', 'ATA', 'ATG',
    'ARG', 'ARM', 'ABW', 'AUS', 'AUT', 'AZE', 'BHS', 'BHR', 'BGD', 'BRB',
    'BLR', 'BEL', 'BLZ', 'BEN', 'BMU', 'BTN', 'BOL', 'BES', 'BIH', 'BWA',
    'BVT', 'BRA', 'IOT', 'BRN', 'BGR', 'BFA', 'BDI', 'CPV', 'KHM', 'CMR',
    'CAN', 'CYM', 'CAF', 'TCD', 'CHL', 'CHN', 'CXR', 'CCK', 'COL', 'COM',
    'COG', 'COD', 'COK', 'CRI', 'CIV', 'HRV', 'CUB', 'CUW', 'CYP', 'CZE',
    'DNK', 'DJI', 'DMA', 'DOM', 'ECU', 'EGY', 'SLV', 'GNQ', 'ERI', 'EST',
    'SWZ', 'ETH', 'FLK', 'FRO', 'FJI', 'FIN', 'FRA', 'GUF', 'PYF', 'ATF',
    'GAB', 'GMB', 'GEO', 'DEU', 'GHA', 'GIB', 'GRC', 'GRL', 'GRD', 'GLP',
    'GUM', 'GTM', 'GGY', 'GIN', 'GNB', 'GUY', 'HTI', 'HMD', 'VAT', 'HND',
    'HKG', 'HUN', 'ISL', 'IND', 'IDN', 'IRN', 'IRQ', 'IRL', 'IMN', 'ISR',
    'ITA', 'JAM', 'JPN', 'JEY', 'JOR', 'KAZ', 'KEN', 'KIR', 'PRK', 'KOR',
    'KWT', 'KGZ', 'LAO', 'LVA', 'LBN', 'LSO', 'LBR', 'LBY', 'LIE', 'LTU',
    'LUX', 'MAC', 'MDG', 'MWI', 'MYS', 'MDV', 'MLI', 'MLT', 'MHL', 'MTQ',
    'MRT', 'MUS', 'MYT', 'MEX', 'FSM', 'MDA', 'MCO', 'MNG', 'MNE', 'MSR',
    'MAR', 'MOZ', 'MMR', 'NAM', 'NRU', 'NPL', 'NLD', 'NCL', 'NZL', 'NIC',
    'NER', 'NGA', 'NIU', 'NFK', 'MKD', 'MNP', 'NOR', 'OMN', 'PAK', 'PLW',
    'PSE', 'PAN', 'PNG', 'PRY', 'PER', 'PHL', 'PCN', 'POL', 'PRT', 'PRI',
    'QAT', 'REU', 'ROU', 'RUS', 'RWA', 'BLM', 'SHN', 'KNA', 'LCA', 'MAF',
    'SPM', 'VCT', 'WSM', 'SMR', 'STP', 'SAU', 'SEN', 'SRB', 'SYC', 'SLE',
    'SGP', 'SXM', 'SVK', 'SVN', 'SLB', 'SOM', 'ZAF', 'SGS', 'SSD', 'ESP',
    'LKA', 'SDN', 'SUR', 'SJM', 'SWE', 'CHE', 'SYR', 'TWN', 'TJK', 'TZA',
    'THA', 'TLS', 'TGO', 'TKL', 'TON', 'TTO', 'TUN', 'TUR', 'TKM', 'TCA',
    'TUV', 'UGA', 'UKR', 'ARE', 'GBR', 'USA', 'UMI', 'URY', 'UZB', 'VUT',
    'VEN', 'VNM', 'VGB', 'VIR', 'WLF', 'ESH', 'YEM', 'ZMB', 'ZWE'
);

-- CREATE TYPE shared.country_code AS ENUM (
--     'AF', 'AX', 'AL', 'DZ', 'AS', 'AD', 'AO', 'AI', 'AQ', 'AG',
--     'AR', 'AM', 'AW', 'AU', 'AT', 'AZ', 'BS', 'BH', 'BD', 'BB',
--     'BY', 'BE', 'BZ', 'BJ', 'BM', 'BT', 'BO', 'BQ', 'BA', 'BW',
--     'BV', 'BR', 'IO', 'BN', 'BG', 'BF', 'BI', 'CV', 'KH', 'CM',
--     'CA', 'KY', 'CF', 'TD', 'CL', 'CN', 'CX', 'CC', 'CO', 'KM',
--     'CG', 'CD', 'CK', 'CR', 'CI', 'HR', 'CU', 'CW', 'CY', 'CZ',
--     'DK', 'DJ', 'DM', 'DO', 'EC', 'EG', 'SV', 'GQ', 'ER', 'EE',
--     'SZ', 'ET', 'FK', 'FO', 'FJ', 'FI', 'FR', 'GF', 'PF', 'TF',
--     'GA', 'GM', 'GE', 'DE', 'GH', 'GI', 'GR', 'GL', 'GD', 'GP',
--     'GU', 'GT', 'GG', 'GN', 'GW', 'GY', 'HT', 'HM', 'VA', 'HN',
--     'HK', 'HU', 'IS', 'IN', 'ID', 'IR', 'IQ', 'IE', 'IM', 'IL',
--     'IT', 'JM', 'JP', 'JE', 'JO', 'KZ', 'KE', 'KI', 'KP', 'KR',
--     'KW', 'KG', 'LA', 'LV', 'LB', 'LS', 'LR', 'LY', 'LI', 'LT',
--     'LU', 'MO', 'MG', 'MW', 'MY', 'MV', 'ML', 'MT', 'MH', 'MQ',
--     'MR', 'MU', 'YT', 'MX', 'FM', 'MD', 'MC', 'MN', 'ME', 'MS',
--     'MA', 'MZ', 'MM', 'NA', 'NR', 'NP', 'NL', 'NC', 'NZ', 'NI',
--     'NE', 'NG', 'NU', 'NF', 'MK', 'MP', 'NO', 'OM', 'PK', 'PW',
--     'PS', 'PA', 'PG', 'PY', 'PE', 'PH', 'PN', 'PL', 'PT', 'PR',
--     'QA', 'RE', 'RO', 'RU', 'RW', 'BL', 'SH', 'KN', 'LC', 'MF',
--     'PM', 'VC', 'WS', 'SM', 'ST', 'SA', 'SN', 'RS', 'SC', 'SL',
--     'SG', 'SX', 'SK', 'SI', 'SB', 'SO', 'ZA', 'GS', 'SS', 'ES',
--     'LK', 'SD', 'SR', 'SJ', 'SE', 'CH', 'SY', 'TW', 'TJ', 'TZ',
--     'TH', 'TL', 'TG', 'TK', 'TO', 'TT', 'TN', 'TR', 'TM', 'TC',
--     'TV', 'UG', 'UA', 'AE', 'GB', 'US', 'UM', 'UY', 'UZ', 'VU',
--     'VE', 'VN', 'VG', 'VI', 'WF', 'EH', 'YE', 'ZM', 'ZW'
-- );

CREATE DOMAIN shared.weekday AS INTEGER CHECK (VALUE BETWEEN 0 AND 6);

CREATE TYPE shared.order_type AS ENUM ('VEHICLE_RENTAL', 'ICE', 'WATER_PURIFIED','WATER_UNPURIFIED', 'FREIGHT', 'BLOCK', 'COCONUT', 'EGG', 'CHARCOAL');

CREATE TABLE IF NOT EXISTS shared.verification_codes
(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    verification_code INT NOT NULL DEFAULT floor(random() * 900000 + 100000)::int,
    failed_attempts SMALLINT NOT NULL DEFAULT 0,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);CREATE SCHEMA IF NOT EXISTS _user;

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
);CREATE SCHEMA IF NOT EXISTS payment;CREATE SCHEMA IF NOT EXISTS organization;
CREATE TYPE organization.bank_account_ownership AS ENUM ('PERSONAL', 'BUSINESS');
CREATE TYPE organization.banK_account_type AS ENUM ('SAVINGS', 'CHEQUING');
CREATE TYPE organization.preferred_water_unit AS ENUM ('GALLON', 'LITRE');
CREATE TYPE organization.preferred_general_unit AS ENUM ('EACH', 'DOZEN');

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
    is_worker BOOLEAN NOT NULL,
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

CREATE TABLE IF NOT EXISTS organization.timeslots
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL,
    start_hh   SMALLINT NOT NULL CHECK (start_hh BETWEEN 0 AND 23),
    start_mm SMALLINT NOT NULL CHECK (start_mm = 0 OR start_mm = 30),
    end_hh   SMALLINT NOT NULL CHECK (end_hh BETWEEN 0 AND 23),
    end_mm SMALLINT NOT NULL CHECK (end_mm = 0 OR end_mm = 30),
    wkdays shared.weekday[] NOT NULL CHECK ( array_length(wkdays, 1) > 0 ),
    average_job_duration BIGINT NOT NULL,
    order_type shared.order_type NOT NULL,
    UNIQUE(organization_id, order_type),
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS organization.water_sales_config
-- (
--     id BIGSERIAL PRIMARY KEY,
--     organization_id BIGSERIAL NOT NULL UNIQUE,
--     preferred_unit organization.preferred_water_unit NOT NULL DEFAULT 'GALLON',
--     water_price_min FLOAT NOT NULL,
--     water_price_per_km FLOAT NOT NULL,
--     water_price_per_gallon FLOAT,
--     unpurified_water_price_per_gallon FLOAT,
--     bottle_cost FLOAT,
--     barrel_cost FLOAT,
--     is_enabled boolean DEFAULT true,
--     CHECK (
--         purified_water_price_per_gallon IS NOT NULL
--         OR unpurified_water_price_per_gallon IS NOT NULL
--         OR bottle_cost IS NOT NULL
--         OR barrel_cost IS NOT NULL
--     ),
--     CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
-- );


CREATE TABLE IF NOT EXISTS organization.potable_water_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_water_unit NOT NULL DEFAULT 'GALLON',
    water_price_min FLOAT NOT NULL,
    water_price_per_km FLOAT NOT NULL,
    water_price_per_gallon FLOAT,
    bottle_cost FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.irrigation_water_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_water_unit NOT NULL DEFAULT 'GALLON',
    water_price_min FLOAT NOT NULL,
    water_price_per_km FLOAT NOT NULL,
    water_price_per_gallon FLOAT,
    barrel_cost FLOAT,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.block_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_general_unit NOT NULL DEFAULT 'EACH',
    block_price_min FLOAT NOT NULL,
    block_price FLOAT NOT NULL,
    block_price_per_km FLOAT NOT NULL,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.coconut_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_general_unit NOT NULL DEFAULT 'EACH',
    coconut_price_min FLOAT NOT NULL,
    coconut_price FLOAT NOT NULL,
    coconut_price_per_km  FLOAT NOT NULL,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.egg_sales_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    preferred_unit organization.preferred_general_unit NOT NULL DEFAULT 'EACH',
    egg_price_min FLOAT NOT NULL,
    egg_price FLOAT NOT NULL,
    egg_price_per_km FLOAT NOT NULL,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.freight_config
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    freight_price_min FLOAT NOT NULL,
    is_enabled boolean DEFAULT true,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.charcoal_sales_config
(
    -- TODO:
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL UNIQUE,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization.sales_units
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL,
    order_type shared.order_type NOT NULL,
    quantity SMALLINT NOT NULL,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE

);

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
CREATE SCHEMA IF NOT EXISTS logistic;

CREATE TYPE logistic.vehicle_fuel_type AS ENUM ('ELECTRIC', 'PETROL', 'DIESEL', 'USLD', 'LNG', 'LPG');
CREATE TYPE logistic.freight_truck_class AS ENUM ('VAN', 'OPEN_TRUCK_SMALL', 'OPEN_TRUCK_LARGE', 'BOX_TRUCK_SMALL', 'BOX_TRUCK_LARGE');

CREATE TABLE IF NOT EXISTS logistic.vehicles 
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL,
    year SMALLINT NOT NULL,
    make TEXT NOT NULL,
    model TEXT NOT NULL,
    chassis_number TEXT NOT NULL,
    license_plate_number TEXT NOT NULL,
    seating_capacity SMALLINT NOT NULL,
    fuel_type logistic.vehicle_fuel_type NOT NULL,
    is_manual_transmission BOOLEAN NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS logistic.trucks 
(
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGSERIAL NOT NULL,
    -- organization_id BIGSERIAL NOT NULL,
    cost FLOAT NOT NULL DEFAULT 0,
    freight_price_per_km FLOAT NOT NULL,
    has_tipper BOOLEAN NOT NULL CHECK (bed_height_mm IS NULL OR has_tipper = FALSE),
    has_crane BOOLEAN NOT NULL CHECK (bed_height_mm IS NULL OR has_crane = FALSE),
    has_bucket BOOLEAN NOT NULL CHECK (bed_height_mm IS NULL OR has_bucket = FALSE),
    has_freezer BOOLEAN NOT NULL CHECK (bed_height_mm IS NOT NULL OR has_freezer = FALSE),
    has_refrigerator BOOLEAN NOT NULL CHECK (bed_height_mm IS NOT NULL OR has_refrigerator = FALSE),
    bed_length_mm FLOAT NOT NULL CHECK (bed_length_mm > 0),
    bed_width_mm FLOAT NOT NULL CHECK (bed_width_mm > 0),
    bed_height_mm FLOAT,
    license_plate_number TEXT NOT NULL,
    chassis_number TEXT NOT NULL,
    last_location GEOGRAPHY(POINT, 4326) NOT NULL,
    -- class logistic.freight_truck_class NOT NULL,
    class TEXT GENERATED ALWAYS AS (
        CASE 
            WHEN bed_length_mm >= 3048 AND bed_height_mm = null THEN 'OPEN_TRUCK_LARGE'
            WHEN bed_length_mm < 3048 AND bed_height_mm = null THEN 'OPEN_TRUCK_SMALL'
            WHEN bed_length_mm >= 3048 AND bed_height_mm != null THEN 'BOX_TRUCK_LARGE'
            WHEN bed_length_mm < 3048 AND bed_height_mm >= 1500 THEN 'BOX_TRUCK_SMALL'
            -- WHEN bed_length_mm < 3048 AND bed_height_mm < 1500 null THEN 'VAN'
            ELSE 'VAN'
        END
    ) STORED,
    -- This casts the generated text to your ENUM to validate it
    CONSTRAINT class_check CHECK (class::logistic.freight_truck_class IS NOT NULL),
    CONSTRAINT fk_vehicle FOREIGN KEY (vehicle_id) REFERENCES logistic.vehicles(id) ON DELETE CASCADE,
    UNIQUE(vehicle_id)
);
CREATE INDEX idx_logistic_last_location ON logistic.trucks USING GIST(last_location);

-- CREATE TABLE IF NOT EXISTS logistic.truck_order_types
-- (
--     id BIGSERIAL PRIMARY KEY,
--     truck_id BIGSERIAL NOT NULL,
--     order_type _order.order_type NOT NULL,
--     CONSTRAINT fk_truck FOREIGN KEY (truck_id) REFERENCES logistic.trucks(id) ON DELETE CASCADE
-- );

-- Table indicating which drivers are permitted to drive which trucks
CREATE TABLE IF NOT EXISTS logistic.truck_assignation (
    id BIGSERIAL PRIMARY KEY,
    truck_id BIGSERIAL NOT NULL,
    driver_id uuid NOT NULL,
    CONSTRAINT fk_truck FOREIGN KEY (truck_id) REFERENCES logistic.trucks(id) ON DELETE CASCADE,
    CONSTRAINT fk_driver FOREIGN KEY (driver_id) REFERENCES organization.accounts(id) ON DELETE CASCADE,
    CONSTRAINT unique_truck_assignation_pair UNIQUE (truck_id, driver_id)
);

CREATE TABLE IF NOT EXISTS logistic.rental_vehicles (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGSERIAL NOT NULL,
    organization_id BIGSERIAL NOT NULL,
    daily_rental_cost FLOAT NOT NULL,
    driver_cost FLOAT,
    is_driver_mandatory BOOLEAN,
    min_days_rental SMALLINT NOT NULL,
    CONSTRAINT chk_driver_logic CHECK (driver_cost IS NOT NULL OR is_driver_mandatory IS NULL),
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS logistic.vehicle_order_type_allocations (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGSERIAL NOT NULL,
    organization_id BIGSERIAL NOT NULL,
    order_type shared.order_type NOT NULL,
    max_capacity FLOAT NOT NULL,
    CONSTRAINT fk_vehicle FOREIGN KEY (vehicle_id) REFERENCES logistic.vehicles(id) ON DELETE CASCADE,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE,
    UNIQUE(vehicle_id, organization_id, order_type)
);CREATE SCHEMA IF NOT EXISTS _order;

CREATE TYPE _order.order_status AS ENUM ('UNPAID', 'PENDING', 'PROCESSING', 'COMPLETE');
CREATE TYPE _order.disconut_code_type AS ENUM ('FIXED', 'PERCENT');
CREATE TYPE _order.payment_method AS ENUM ('PAYPAL',  'BANK_TRANSFER', 'BITCOIN', 'BITCOIN_LIGHTNING');

-- CREATE TABLE IF NOT EXISTS _order.customers 
-- (
--     id BIGSERIAL PRIMARY KEY,
--     name TEXT NOT NULL,
--     phone TEXT NOT NULL,
--     last_ip_address TEXT NOT NULL,
--     email TEXT,
--     password_hash VARCHAR(128), -- Field is nullable due to Google OAUTH login potential
--     google_user_id TEXT,
--     oauth_refresh_token_encrypted TEXT,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
-- );

CREATE TABLE IF NOT EXISTS _order.discount_codes
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL,
    code CHARACTER(16) UNIQUE NOT NULL,
    discount_code_type _order.disconut_code_type NOT NULL,
    value FLOAT NOT NULL CHECK (
        (discount_code_type = 'PERCENT' AND value > 0 AND value < 1)
        OR
        (discount_code_type = 'FIXED' AND value > 0)
    ),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS _order.orders 
(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGSERIAL NOT NULL,
    organization_id BIGSERIAL NOT NULL,
    payment_date BIGINT DEFAULT NULL,
    payment_method _order.payment_method NOT NULL,
    job_date BIGINT CHECK (job_date > 0), -- job date in milliseconds
    job_month SMALLINT CHECK (job_month BETWEEN 0 AND 11),
    job_year SMALLINT CHECK (job_year >= 2025),
    distance_round_trip_in_km FLOAT NOT NULL,
    type shared.order_type NOT NULL,
    status _order.order_status NOT NULL DEFAULT 'PENDING' CHECK ((status = 'PENDING' AND payment_date IS NULL) OR (status != 'PENDING' AND payment_date IS NOT NULL)),
    -- delivery_address TEXT NOT NULL,
    total_in_usd FLOAT NOT NULL,
    paypal_order_id TEXT UNIQUE,
    bank_account_id BIGSERIAL,
    discount_code_id BIGSERIAL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT chk_payment_method_required CHECK (bank_account_id IS NOT NULL OR paypal_order_id IS NOT NULL),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES _user.users(id) ON DELETE CASCADE,
    CONSTRAINT fk_bank_account FOREIGN KEY (bank_account_id) REFERENCES organization.bank_accounts(id) ON DELETE CASCADE,
    CONSTRAINT fk_discount_code FOREIGN KEY (discount_code_id) REFERENCES _order.discount_codes(id) ON DELETE CASCADE,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS _order.water_orders
-- (
--     id BIGSERIAL PRIMARY KEY,
--     order_id uuid NOT NULL UNIQUE,
--     bottles_required BOOLEAN NOT NULL DEFAULT FALSE CHECK (bottles_required = FALSE OR quantity_of_purified_water > 0),
--     quantity_of_purified_water SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_purified_water >= 0),
--     quantity_of_unpurified_water SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_unpurified_water >= 0),
--     quantity_of_barrels SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_barrels >= 0),
--     quantity_of_bottles SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_bottles >= 0),
--     CONSTRAINT fkorder FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
-- );

CREATE TABLE IF NOT EXISTS _order.potable_water_orders
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    quantity_of_water SMALLINT NOT NULL,
    bottles_required BOOLEAN NOT NULL DEFAULT FALSE CHECK (bottles_required = FALSE OR quantity_of_water > 0),
    CONSTRAINT fkorder FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _order.irrigation_water_orders
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    quantity_of_water SMALLINT NOT NULL CHECK (quantity_of_water > 0),
    barrels_required BOOLEAN NOT NULL DEFAULT FALSE CHECK (barrels_required = FALSE OR quantity_of_water > 0),
    quantity_of_barrels SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_barrels >= 0),
    CONSTRAINT fkorder FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _order.ice_orders
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    quantity_of_ice SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_ice >= 0),
    CONSTRAINT fkorder FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _order.rental_vehicle_orders
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    rental_vehicle_id BIGSERIAL NOT NULL,
    start_date BIGINT CHECK (start_date > 0),
    end_date BIGINT CHECK (end_date > 0),
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_rental_vehicle FOREIGN KEY (rental_vehicle_id) REFERENCES logistic.rental_vehicles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _order.order_waypoints
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL,
    waypoint GEOGRAPHY(POINT, 4326) NOT NULL,
    address TEXT NOT NULL,
    -- TODO: Add human readable address ?
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
);
CREATE INDEX idx_order_waypoints_geog ON _order.order_waypoints USING GIST(waypoint);

CREATE TABLE IF NOT EXISTS _order.freight_orders
(
    id BIGSERIAL PRIMARY KEY,
    truck_id BIGSERIAL NOT NULL,
    order_id uuid NOT NULL UNIQUE,
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_truck FOREIGN KEY (truck_id) REFERENCES logistic.trucks(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _order.block_orders
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    quantity_of_blocks SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_blocks > 0),
    CONSTRAINT fkorder FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _order.coconut_orders
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    quantity_of_coconuts SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_coconuts > 0),
    CONSTRAINT fkorder FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _order.egg_orders
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    quantity_of_eggs SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_eggs > 0),
    CONSTRAINT fkorder FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _order.charcoal_orders
(
    id BIGSERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    quantity_of_charcoal SMALLINT NOT NULL DEFAULT 0 CHECK (quantity_of_charcoal > 0),
    CONSTRAINT fkorder FOREIGN KEY (order_id) REFERENCES _order.orders(id) ON DELETE CASCADE
);