CREATE SCHEMA IF NOT EXISTS _order;

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
    comment VARCHAR(100) DEFAULT '',
    rating SMALLINT CHECK (rating BETWEEN 1 AND 5),
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

CREATE TABLE IF NOT EXISTS _order.enrollment
(
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGSERIAL NOT NULL UNIQUE,
    freight BOOLEAN NOT NULL DEFAULT FALSE,
    water_potable BOOLEAN NOT NULL DEFAULT FALSE,
    water_unpurified BOOLEAN NOT NULL DEFAULT FALSE,
    ice BOOLEAN NOT NULL DEFAULT FALSE,
    blocks BOOLEAN NOT NULL DEFAULT FALSE,
    coconut BOOLEAN NOT NULL DEFAULT FALSE,
    egg BOOLEAN NOT NULL DEFAULT FALSE,
    charcoal BOOLEAN NOT NULL DEFAULT FALSE,
    cooking_gas BOOLEAN NOT NULL DEFAULT FALSE,
    honey BOOLEAN NOT NULL DEFAULT FALSE,
    cesspool BOOLEAN NOT NULL DEFAULT FALSE,
    construction BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_vehicle_id FOREIGN KEY (vehicle_id) REFERENCES logistic.vehicles (id) ON DELETE CASCADE
);