CREATE SCHEMA IF NOT EXISTS logistic;

CREATE TYPE logistic.vehicle_fuel_type AS ENUM ('ELECTRIC', 'PETROL', 'DIESEL', 'USLD', 'LNG', 'LPG');
CREATE TYPE logistic.freight_truck_class AS ENUM ('VAN', 'OPEN_TRUCK_SMALL', 'OPEN_TRUCK_LARGE', 'BOX_TRUCK_SMALL', 'BOX_TRUCK_LARGE');
-- CREATE TYPE logistic.freight_truck_property AS ENUM ('NONE', 'TIPPER', 'FREEZER', 'REFRIGERATED', 'BUCKET', 'CRANE', 'CESSPOOL');

CREATE TABLE IF NOT EXISTS logistic.vehicles 
(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL,
    driver_account_id uuid NOT NULL,
    year SMALLINT NOT NULL,
    make TEXT NOT NULL,
    model TEXT NOT NULL,
    chassis_number TEXT NOT NULL,
    license_plate_number TEXT NOT NULL,
    seating_capacity SMALLINT NOT NULL,
    fuel_type logistic.vehicle_fuel_type NOT NULL,
    is_manual_transmission BOOLEAN NOT NULL,
    is_approved BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization.organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_driver FOREIGN KEY (driver_account_id) REFERENCES organization.accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS logistic.trucks 
(
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGSERIAL NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    cost FLOAT NOT NULL DEFAULT 0,
    freight_price_per_km FLOAT NOT NULL,
    property BIGSERIAL NOT NULL,
    bed_length_mm FLOAT NOT NULL CHECK (bed_length_mm > 0),
    bed_width_mm FLOAT NOT NULL CHECK (bed_width_mm > 0),
    bed_height_mm FLOAT,
    last_location GEOGRAPHY(POINT, 4326) NOT NULL,
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
    CONSTRAINT fk_property FOREIGN KEY (property) REFERENCES logistic.truck_properties(id) ON DELETE CASCADE,
    UNIQUE(vehicle_id)
);
CREATE INDEX idx_logistic_last_location ON logistic.trucks USING GIST(last_location);

CREATE TABLE IF NOT EXISTS logistic.truck_properties (
    id BIGSERIAL PRIMARY KEY,
    property TEXT NOT NULL
);

-- Get


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
    daily_rental_cost FLOAT NOT NULL,
    driver_cost FLOAT,
    is_driver_mandatory BOOLEAN,
    min_days_rental SMALLINT NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT chk_driver_logic CHECK (driver_cost IS NOT NULL OR is_driver_mandatory IS NULL),
    UNIQUE(vehicle_id)
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
);

CREATE TABLE IF NOT EXISTS logistic.timeslots
(
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGSERIAL NOT NULL,
    start_hh   SMALLINT NOT NULL CHECK (start_hh BETWEEN 0 AND 23),
    start_mm SMALLINT NOT NULL CHECK (start_mm = 0 OR start_mm = 30),
    end_hh   SMALLINT NOT NULL CHECK (end_hh BETWEEN 0 AND 23),
    end_mm SMALLINT NOT NULL CHECK (end_mm = 0 OR end_mm = 30),
    wkdays shared.weekday[] NOT NULL CHECK ( array_length(wkdays, 1) > 0 ),
    CONSTRAINT fk_vehicle FOREIGN KEY (vehicle_id) REFERENCES logistic.vehicles(id) ON DELETE CASCADE
);