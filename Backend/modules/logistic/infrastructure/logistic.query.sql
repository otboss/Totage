-- name: CreateVehicle :one
INSERT INTO logistic.vehicles (
    organization_id,
    driver_account_id,
    year,
    make,
    model,
    chassis_number,
    license_plate_number,
    seating_capacity,
    fuel_type,
    is_manual_transmission
) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
)
RETURNING 
    id;

-- name: UpdateVehicleDriver :exec
UPDATE logistic.vehicles 
SET driver_account_id = $1
WHERE id = sqlc.arg(vehicle_id)::bigint;

-- name: GetTruck :one
SELECT 
    vehicle_id,
    cost,
    freight_price_per_km,
    property,
    bed_length_mm,
    bed_width_mm,
    bed_height_mm,
    is_public,
    ST_X(last_location::geometry)::float AS longitude,
    ST_Y(last_location::geometry)::float AS latitude
FROM logistic.trucks WHERE sqlc.arg(truck_id)::bigint;

-- name: GetTruckByVehicleId :one
SELECT 
    id,
    vehicle_id,
    cost,
    freight_price_per_km,
    property,
    bed_length_mm,
    bed_width_mm,
    bed_height_mm,
    is_public,
    ST_X(last_location::geometry)::float AS longitude,
    ST_Y(last_location::geometry)::float AS latitude
FROM logistic.trucks 
WHERE vehicle_id = sqlc.arg(vehicle_id)::bigint;


-- name: GetTruckPropertyByProperty :one
SELECT
    id,
    property
FROM logistic.truck_properties
WHERE property = $1;


-- name: GetTruckPropertyByPropertyId :one
SELECT
    id,
    property
FROM logistic.truck_properties
WHERE id = $1;


-- name: GetTruckCount :one
SELECT
    COUNT(*)
FROM logistic.trucks 
JOIN logistic.vehicles
ON logistic.trucks.vehicle_id = logistic.vehicles.id
WHERE logistic.vehicles.organization_id = $1;

-- -- name: GetTruckByBedDimensions :one
-- SELECT 
--     id,
--     cost,
--     freight_price_per_km,
--     
--     has_tipper,
--     has_crane,
--     has_bucket,
--     has_freezer,
--     has_refrigerator,
--     bed_length_mm,
--     bed_width_mm,
--     bed_height_mm,
--     logistic.trucks.chassis_number,
--     license_plate_number
-- FROM logistic.trucks
-- WHERE organization_id = $1 
-- AND has_tipper = $2 
-- AND has_crane = $3  
-- AND has_bucket = $4 
-- AND has_freezer = $5 
-- AND has_refrigerator = $6
-- AND bed_length_mm = $7 
-- AND bed_width_mm = $8
-- AND (
--     (sqlc.arg(bed_height_mm)::float IS NULL AND bed_height_mm IS NULL) OR
--     bed_height_mm = sqlc.arg(bed_height_mm)::float
-- );

-- name: GetTrucksByOrganization :many
SELECT
    logistic.trucks.id,
    logistic.trucks.vehicle_id,
    logistic.trucks.cost,
    logistic.trucks.freight_price_per_km,
    logistic.trucks.property,
    logistic.trucks.bed_length_mm,
    logistic.trucks.bed_width_mm,
    logistic.trucks.bed_height_mm,
    logistic.trucks.is_public,
    logistic.trucks.ST_X(last_location::geometry)::float AS longitude,
    logistic.trucks.ST_Y(last_location::geometry)::float AS latitude
FROM logistic.trucks
JOIN logistic.vehicles
ON logistic.trucks.vehicle_id = logistic.vehicles.id
WHERE logistic.vehicles.organization_id = $1
AND logistic.trucks.id >= sqlc.arg(offset_id)::bigint
-- ORDER BY logistic.trucks.bed_length_mm ASC
LIMIT logistic.trucks.id >= sqlc.arg(limitation)::int;

-- name: GetTruckClassesNearby :many
SELECT
    class,
    property
FROM logistic.trucks
WHERE ST_DWithin(
    last_location, 
    ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY, 
    100000
)
GROUP BY class
ORDER BY MIN(last_location <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY)
LIMIT (
    SELECT array_length(enum_range(NULL::logistic.freight_truck_class), 1)::int
);

-- SELECT DISTINCT
--     class
-- FROM logistic.trucks
-- WHERE ST_DWithin(
--     address_geog, 
--     ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY, 
--     100000
-- )
-- ORDER BY address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
-- LIMIT (
--     SELECT array_length(enum_range(NULL::logistic.freight_truck_class), 1)::int
-- );

-- name: GetTruckByOrganizationAndVehicleId :one
SELECT
    logistic.trucks.id,
    logistic.trucks.vehicle_id,
    logistic.trucks.cost,
    logistic.trucks.freight_price_per_km,
    logistic.trucks.property,
    logistic.trucks.bed_length_mm,
    logistic.trucks.bed_width_mm,
    logistic.trucks.bed_height_mm,
    logistic.trucks.is_public,
    logistic.trucks.ST_X(last_location::geometry)::float AS longitude,
    logistic.trucks.ST_Y(last_location::geometry)::float AS latitude
FROM logistic.trucks
JOIN logistic.vehicles
ON logistic.trucks.vehicle_id = logistic.vehicles.id
WHERE logistic.vehicles.organization_id = $1
AND logistic.trucks.vehicle_id = $2;

-- name: GetVehicle :one
SELECT
    id,
    organization_id,
    driver_account_id,
    year,
    make,
    model,
    chassis_number,
    license_plate_number,
    seating_capacity,
    fuel_type,
    license_plate_number,
    chassis_number,
    is_manual_transmission,
    is_approved
FROM logistic.vehicles
WHERE id = $1;

-- name: GetVehicles :many
SELECT
    id,
    organization_id,
    driver_account_id,
    year,
    make,
    model,
    chassis_number,
    license_plate_number,
    seating_capacity,
    fuel_type,
    license_plate_number,
    chassis_number,
    is_manual_transmission,
    is_approved
FROM logistic.vehicles
WHERE organization_id = $1
AND id >= sqlc.arg(offset_id)::bigint
AND (
    logistic.vehicles.make ILIKE '%' || sqlc.arg(filter_text)::text || '%'
    OR
    logistic.vehicles.model ILIKE '%' || sqlc.arg(filter_text)::text || '%'
)
-- ORDER BY logistic.rental_vehicles.id ASC
LIMIT sqlc.arg(limitation)::int;


-- -- name: GetTrucksFreightClosestToPoint :many
-- SELECT
--     logistic.trucks.id,
--     logistic.trucks.vehicle_id,
--     logistic.trucks.cost,
--     logistic.trucks.freight_price_per_km,
--     logistic.trucks.has_tipper,
--     logistic.trucks.has_crane,
--     logistic.trucks.has_bucket,
--     logistic.trucks.has_freezer,
--     logistic.trucks.has_refrigerator,
--     logistic.trucks.bed_length_mm,
--     logistic.trucks.bed_width_mm,
--     logistic.trucks.bed_height_mm,
--     logistic.trucks.license_plate_number,
--     logistic.trucks.class,
--     ST_Distance(
--         organization.organizations.address_geog,
--         ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
--     )::float AS distance_meters
-- FROM logistic.trucks
-- JOIN organization.organizations
-- ON logistic.trucks.organization_id = organization.organizations.id
-- JOIN logistic.vehicle_order_type_allocations
-- ON logistic.vehicle_order_type_allocations.vehicle_id = logistic.trucks.vehicle_id
-- AND logistic.vehicle_order_type_allocations.organization_id = logistic.trucks.organization_id
-- WHERE organization.organizations.is_freight_service_supported
-- AND organization.organizations.is_approved
-- AND (
--     organization.organizations.max_order_queue_size == NULL OR
--     (SELECT COUNT(*) FROM _order.orders WHERE _order.orders.status != 'COMPLETE') < organization.organizations.max_order_queue_size
-- )
-- AND logistic.vehicle_order_type_allocations.order_type = 'FREIGHT'::shared.order_type
-- AND logistic.vehicle_order_type_allocations.max_capacity <= sqlc.arg(max_capacity)::float
-- ORDER BY address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
-- LIMIT sqlc.arg(limitation)::int;


-- -- name: GetTrucksEggsClosestToPoint :many
-- SELECT
--     logistic.trucks.id,
--     logistic.trucks.vehicle_id,
--     logistic.trucks.cost,
--     logistic.trucks.freight_price_per_km,
    
--     logistic.trucks.has_tipper,
--     logistic.trucks.has_crane,
--     logistic.trucks.has_bucket,
--     logistic.trucks.has_freezer,
--     logistic.trucks.has_refrigerator,
--     logistic.trucks.bed_length_mm,
--     logistic.trucks.bed_width_mm,
--     logistic.trucks.bed_height_mm,
--     logistic.trucks.license_plate_number,
--     ST_Distance(
--         organization.organizations.address_geog,
--         ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
--     )::float AS distance_meters
-- FROM logistic.trucks
-- JOIN organization.organizations
-- ON logistic.trucks.organization_id = organization.organizations.id
-- WHERE organization.organizations.is_egg_service_supported 
-- AND organization.organizations.is_approved
-- AND (
--     organization.organizations.max_order_queue_size == NULL OR
--     (SELECT COUNT(*) FROM _order.orders WHERE _order.orders.status != 'COMPLETE') < organization.organizations.max_order_queue_size
-- )
-- AND logistic.vehicle_order_type_allocations.order_type = 'EGG'::shared.order_type
-- AND logistic.vehicle_order_type_allocations.max_capacity <= sqlc.arg(max_capacity)::float
-- ORDER BY address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
-- LIMIT sqlc.arg(limitation)::int;

-- -- name: GetTrucksWaterPurifiedClosestToPoint :many
-- SELECT
--     logistic.trucks.id,
--     logistic.trucks.vehicle_id,
--     logistic.trucks.cost,
--     logistic.trucks.freight_price_per_km,
    
--     logistic.trucks.has_tipper,
--     logistic.trucks.has_crane,
--     logistic.trucks.has_bucket,
--     logistic.trucks.has_freezer,
--     logistic.trucks.has_refrigerator,
--     logistic.trucks.bed_length_mm,
--     logistic.trucks.bed_width_mm,
--     logistic.trucks.bed_height_mm,
--     logistic.trucks.license_plate_number,
--     logistic.trucks.class,
--     ST_Distance(
--         organization.organizations.address_geog,
--         ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
--     )::float AS distance_meters
-- FROM logistic.trucks
-- JOIN organization.organizations
-- ON logistic.trucks.organization_id = organization.organizations.id
-- JOIN logistic.vehicle_order_type_allocations
-- ON logistic.vehicle_order_type_allocations.vehicle_id = logistic.trucks.vehicle_id
-- AND logistic.vehicle_order_type_allocations.organization_id = logistic.trucks.organization_id
-- WHERE organization.organizations.is_water_service_supported
-- AND organization.organizations.is_approved
-- AND (
--     organization.organizations.max_order_queue_size == NULL OR
--     (SELECT COUNT(*) FROM _order.orders WHERE _order.orders.status != 'COMPLETE') < organization.organizations.max_order_queue_size
-- )
-- AND logistic.vehicle_order_type_allocations.order_type = 'WATER_PURIFIED'::shared.order_type
-- AND logistic.vehicle_order_type_allocations.max_capacity <= sqlc.arg(max_capacity)::float
-- ORDER BY address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
-- LIMIT sqlc.arg(limitation)::int;

-- -- name: GetTrucksWaterUnpurifiedClosestToPoint :many
-- SELECT
--     logistic.trucks.id,
--     logistic.trucks.vehicle_id,
--     logistic.trucks.cost,
--     logistic.trucks.freight_price_per_km,
--     logistic.trucks.has_tipper,
--     logistic.trucks.has_crane,
--     logistic.trucks.has_bucket,
--     logistic.trucks.has_freezer,
--     logistic.trucks.has_refrigerator,
--     logistic.trucks.bed_length_mm,
--     logistic.trucks.bed_width_mm,
--     logistic.trucks.bed_height_mm,
--     logistic.trucks.license_plate_number,
--     logistic.trucks.class,
--     ST_Distance(
--         organization.organizations.address_geog,
--         ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
--     )::float AS distance_meters
-- FROM logistic.trucks
-- JOIN organization.organizations
-- ON logistic.trucks.organization_id = organization.organizations.id
-- JOIN logistic.vehicle_order_type_allocations
-- ON logistic.vehicle_order_type_allocations.vehicle_id = logistic.trucks.vehicle_id
-- AND logistic.vehicle_order_type_allocations.organization_id = logistic.trucks.organization_id
-- WHERE organization.organizations.is_water_service_supported
-- AND organization.organizations.is_approved
-- AND (
--     organization.organizations.max_order_queue_size == NULL OR
--     (SELECT COUNT(*) FROM _order.orders WHERE _order.orders.status != 'COMPLETE') < organization.organizations.max_order_queue_size
-- )
-- AND logistic.vehicle_order_type_allocations.order_type = 'WATER_UNPURIFIED'::shared.order_type
-- AND logistic.vehicle_order_type_allocations.max_capacity <= sqlc.arg(max_capacity)::float
-- ORDER BY address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
-- LIMIT sqlc.arg(limitation)::int;

-- -- name: GetTrucksBlocksClosestToPoint :many
-- SELECT
--     logistic.trucks.id,
--     logistic.trucks.vehicle_id,
--     logistic.trucks.cost,
--     logistic.trucks.freight_price_per_km,
--     logistic.trucks.has_tipper,
--     logistic.trucks.has_crane,
--     logistic.trucks.has_bucket,
--     logistic.trucks.has_freezer,
--     logistic.trucks.has_refrigerator,
--     logistic.trucks.bed_length_mm,
--     logistic.trucks.bed_width_mm,
--     logistic.trucks.bed_height_mm,
--     logistic.trucks.license_plate_number,
--     logistic.trucks.class,
--     ST_Distance(
--         organization.organizations.address_geog,
--         ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
--     )::float AS distance_meters
-- FROM logistic.trucks
-- JOIN organization.organizations
-- ON logistic.trucks.organization_id = organization.organizations.id
-- JOIN logistic.vehicle_order_type_allocations
-- ON logistic.vehicle_order_type_allocations.vehicle_id = logistic.trucks.vehicle_id
-- AND logistic.vehicle_order_type_allocations.organization_id = logistic.trucks.organization_id
-- WHERE organization.organizations.is_block_service_supported
-- AND organization.organizations.is_approved
-- AND (
--     organization.organizations.max_order_queue_size == NULL OR
--     (SELECT COUNT(*) FROM _order.orders WHERE _order.orders.status != 'COMPLETE') < organization.organizations.max_order_queue_size
-- )
-- AND logistic.vehicle_order_type_allocations.order_type = 'BLOCK'::shared.order_type
-- AND logistic.vehicle_order_type_allocations.max_capacity <= sqlc.arg(max_capacity)::float
-- ORDER BY address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
-- LIMIT sqlc.arg(limitation)::int;

-- name: GetTrucksClosestToPoint :many
SELECT
    logistic.trucks.id,
    logistic.trucks.vehicle_id,
    logistic.trucks.cost,
    logistic.trucks.freight_price_per_km,
    logistic.trucks.property,
    logistic.trucks.bed_length_mm,
    logistic.trucks.bed_width_mm,
    logistic.trucks.bed_height_mm,
    logistic.trucks.is_public,
    logistic.trucks.class,
    ST_Distance(
        organization.organizations.address_geog,
        ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
    )::float AS distance_meters
FROM logistic.trucks
JOIN organization.organizations
ON logistic.trucks.organization_id = organization.organizations.id
JOIN logistic.vehicle_order_type_allocations
ON logistic.vehicle_order_type_allocations.vehicle_id = logistic.trucks.vehicle_id
AND logistic.vehicle_order_type_allocations.organization_id = logistic.trucks.organization_id
WHERE organization.organizations.is_approved
AND (
    organization.organizations.max_order_queue_size == NULL OR
    (SELECT COUNT(*) FROM _order.orders WHERE _order.orders.status != 'COMPLETE') < organization.organizations.max_order_queue_size
)
AND logistic.vehicle_order_type_allocations.order_type = sqlc.arg(order_type)::shared.order_type
AND logistic.vehicle_order_type_allocations.max_capacity <= sqlc.arg(max_capacity)::float
ORDER BY address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
LIMIT sqlc.arg(limitation)::int;

-- -- name: GetTrucksBlocksClosestToPoint :many
-- SELECT
--     logistic.trucks.id,
--     logistic.trucks.vehicle_id,
--     logistic.trucks.cost,
--     logistic.trucks.freight_price_per_km,
--     logistic.trucks.has_tipper,
--     logistic.trucks.has_crane,
--     logistic.trucks.has_bucket,
--     logistic.trucks.has_freezer,
--     logistic.trucks.has_refrigerator,
--     logistic.trucks.bed_length_mm,
--     logistic.trucks.bed_width_mm,
--     logistic.trucks.bed_height_mm,
--     logistic.trucks.license_plate_number,
--     ST_Distance(
--         organization.organizations.address_geog,
--         ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
--     )::float AS distance_meters
-- FROM logistic.trucks
-- JOIN organization.organizations
-- ON logistic.trucks.organization_id = organization.organizations.id
-- WHERE organization.organizations.is_block_service_supported
-- AND organization.organizations.is_approved
-- AND (
--     organization.organizations.max_order_queue_size == NULL OR
--     (SELECT COUNT(*) FROM _order.orders WHERE _order.orders.status != 'COMPLETE') < organization.organizations.max_order_queue_size
-- )
-- ORDER BY organization.organizations.address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
-- LIMIT sqlc.arg(limitation)::int;


-- name: GetRentalVehiclesPublicClosestToPoint :many
SELECT
    logistic.rental_vehicles.id,
    logistic.rental_vehicles.vehicle_id,
    logistic.vehicles.year,
    logistic.vehicles.make,
    logistic.vehicles.model,
    -- logistic.rental_vehicles.seating_capacity,
    logistic.vehicles.fuel_type,
    logistic.rental_vehicles.daily_rental_cost,
    logistic.rental_vehicles.driver_cost,
    logistic.rental_vehicles.min_days_rental,
    logistic.rental_vehicles.is_driver_mandatory,
    logistic.rental_vehicles.is_public,
    logistic.vehicles.license_plate_number,
    logistic.vehicles.is_manual_transmission,
    logistic.vehicle_order_type_allocations.max_capacity,
    ST_Distance(
        organization.organizations.address_geog,
        ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
    )::float AS distance_meters
FROM logistic.rental_vehicles
JOIN organization.organizations
ON logistic.rental_vehicles.organization_id = organization.organizations.id
JOIN logistic.vehicle_order_type_allocations
ON logistic.vehicle_order_type_allocations.vehicle_id = logistic.trucks.vehicle_id
JOIN logistic.vehicles
ON logistic.rental_vehicles.vehicle_id = logistic.vehicles.id
WHERE organization.organizations.is_rental_vehicle_service_supported 
AND logistic.vehicle_order_type_allocations.organization_id = logistic.trucks.organization_id
AND organization.organizations.is_approved
AND (
    organization.organizations.max_order_queue_size == NULL OR
    (SELECT COUNT(*) FROM _order.orders WHERE _order.orders.status != 'COMPLETE') < organization.organizations.max_order_queue_size
)
AND logistic.vehicle_order_type_allocations.order_type = 'VEHICLE_RENTAL'::shared.order_type
AND logistic.vehicle_order_type_allocations.max_capacity >= sqlc.arg(seating_capacity)::bigint
ORDER BY organization.organizations.address_geog <-> ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float)::GEOGRAPHY
LIMIT sqlc.arg(limitation)::int;

-- name: CreateTruck :one
INSERT INTO logistic.trucks
(
    vehicle_id,
    cost,
    freight_price_per_km,
    property,
    bed_length_mm,
    bed_width_mm,
    bed_height_mm,
    is_public,
    last_location
)
VALUES
(
    $1, $2, $3, $4, $5, $6, $7, $8, ST_SetSRID(ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float), 4326)::geography
)
ON CONFLICT (vehicle_id)
DO UPDATE SET
    cost = $2,
    freight_price_per_km = $3,
    is_public = $8,
    last_location = ST_SetSRID(ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float), 4326)::geography
RETURNING 
    id;

-- name: CreateRentalVehicle :one
INSERT INTO logistic.rental_vehicles
(
    vehicle_id,
    daily_rental_cost,
    driver_cost,
    is_driver_mandatory,
    min_days_rental,
    is_public
)
VALUES
(
    $1, $2, $3, $4, $5, $6
)
RETURNING 
    id;

-- name: UpdateTruckLocation :exec
UPDATE logistic.trucks
SET last_location = ST_SetSRID(ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float), 4326)::geography
WHERE id = $1;

-- name: ToggleTruckPublicity :exec
UPDATE logistic.trucks
SET is_public = $1
WHERE vehicle_id = $2;

-- name: ToggleRentalVehiclePublicity :exec
UPDATE logistic.rental_vehicles
SET is_public = $1
WHERE vehicle_id = $2;

-- name: GetRentalVehiclesByOrganizationId :many
SELECT
    logistic.rental_vehicles.id,
    logistic.rental_vehicles.vehicle_id,
    logistic.vehicles.year,
    logistic.vehicles.make,
    logistic.vehicles.model,
    logistic.vehicles.fuel_type,
    logistic.rental_vehicles.daily_rental_cost,
    logistic.rental_vehicles.driver_cost,
    logistic.rental_vehicles.min_days_rental,
    logistic.vehicles.license_plate_number,
    logistic.vehicles.chassis_number,
    logistic.rental_vehicles.is_driver_mandatory,
    logistic.rental_vehicles.is_public,
    logistic.vehicles.is_manual_transmission
FROM logistic.rental_vehicles
JOIN logistic.vehicles
ON logistic.rental_vehicles.vehicle_id = logistic.vehicles.id
WHERE logistic.vehicles.organization_id = $1
AND logistic.rental_vehicles.id >= sqlc.arg(offset_id)::bigint
-- ORDER BY logistic.rental_vehicles.id ASC
LIMIT logistic.rental_vehicles.id >= sqlc.arg(limitation)::int;

-- name: GetRentalVehicle :one
SELECT
    logistic.rental_vehicles.id,
    logistic.rental_vehicles.vehicle_id,
    logistic.vehicles.year,
    logistic.vehicles.make,
    logistic.vehicles.model,
    logistic.vehicles.fuel_type,
    logistic.rental_vehicles.daily_rental_cost,
    logistic.rental_vehicles.driver_cost,
    logistic.rental_vehicles.min_days_rental,
    logistic.vehicles.license_plate_number,
    logistic.vehicles.chassis_number,
    logistic.rental_vehicles.is_driver_mandatory,
    logistic.rental_vehicles.is_public,
    logistic.vehicles.is_manual_transmission
FROM logistic.rental_vehicles
JOIN logistic.vehicles
ON logistic.rental_vehicles.vehicle_id = logistic.vehicles.id
WHERE logistic.rental_vehicles.id = $1;

-- name: GetRentalVehicleByVehicleId :one
SELECT
    id,
    vehicle_id,
    daily_rental_cost,
    driver_cost,
    is_driver_mandatory,
    min_days_rental,
    is_public
FROM logistic.rental_vehicles
WHERE vehicle_id = sqlc.arg(vehicle_id)::bigint;

-- name: GetRentalVehiclesEnabled :one
SELECT EXISTS (
    SELECT 1
    FROM logistic.rental_vehicles
    JOIN logistic.vehicles
    ON logistic.rental_vehicles.vehicle_id = logistic.vehicles.id
    WHERE logistic.vehicles.organization_id = $1
    AND logistic.rental_vehicles.is_public = TRUE
)::bool AS is_enabled;

-- name: GetVehicleOrderTypeAllocations :many
SELECT
    id,
    vehicle_id,
    organization_id,
    max_capacity
FROM (
    SELECT DISTINCT ON (organization_id)
    id,
    vehicle_id,
    organization_id,
    max_capacity
    FROM logistic.vehicle_order_type_allocations
) AS unique_rows
WHERE order_type = $1
AND max_capacity <= $2
ORDER BY RANDOM()
LIMIT 10;

-- name: GetVehicleOrderTypeAllocation :one
SELECT
    id,
    max_capacity
FROM logistic.vehicle_order_type_allocations
WHERE vehicle_id = $1
AND organization_id = $2
AND order_type = $3
LIMIT 1;

-- name: CreateVehicleOrderTypeAllocation :one
INSERT INTO logistic.vehicle_order_type_allocations (
    vehicle_id,
    organization_id,
    order_type,
    max_capacity
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING 
    id;

-- name: DeleteTimeslots :exec
DELETE FROM logistic.timeslots
WHERE vehicle_id = $1;

-- name: DeleteTimeslot :exec
DELETE FROM logistic.timeslots
WHERE id = $1;

-- name: CreateTimeslot :one
INSERT INTO logistic.timeslots
(
    vehicle_id,
    start_hh,
    start_mm,
    end_hh,
    end_mm,
    wkdays
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    sqlc.arg(wkdays)::int[]
)
RETURNING 
    id;

-- name: GetTimeslots :many
SELECT
    id,
    start_hh,
    start_mm,
    end_hh,
    end_mm,
    (wkdays::int[]) as weekdays
FROM logistic.timeslots
WHERE vehicle_id = $1;

-- name: GetTimeslot :one
SELECT
    id,
    start_hh,
    start_mm,
    end_hh,
    end_mm,
    (wkdays::int[]) as weekdays
FROM logistic.timeslots
WHERE id = $1;

-- name: DeleteVehicle :exec
DELETE FROM logistic.vehicles
WHERE id = $1;