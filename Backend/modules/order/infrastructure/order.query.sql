

-- -- name: CreateCustomer :one
-- INSERT INTO _order.customers
--     (
--     name, phone, email, last_ip_address
--     )
-- VALUES
--     (
--         $1, $2, $3, $4
--     )
-- RETURNING
--     id;

-- -- name: CreateCustomer :one
-- INSERT INTO _order.customers
--     (
--         name, phone, last_ip_address, email, password_hash, google_user_id, oauth_refresh_token_encrypted
--     )
-- VALUES
--     (
--         $1, $2, $3, $4, $5, $6, $7 
--     )
-- RETURNING 
--     id;

-- -- name: CreateCustomerWithTimeBuffer :one
-- INSERT INTO _order.customers
--     (
--     name, phone, email, last_ip_address
--     )
-- SELECT $1, $2, $3, $4
-- WHERE NOT EXISTS (
--   SELECT 1
--   FROM _order.customers
--   WHERE last_ip_address = $4 AND created_at > NOW() - INTERVAL '5 seconds'
-- )
-- RETURNING 
--     id;

-- -- name: GetCustomer :one
-- SELECT name, phone, email, last_ip_address FROM _order.customers WHERE id = $1;

-- name: GetIncompleteOrderCountByOrganization :one
SELECT COUNT(*)::int as total_incomplete_orders FROM _order.orders WHERE organization_id = $1 AND status != 'COMPLETE';

-- -- name: GetRecentCustomerByIpAddress :one
-- SELECT name, phone, email FROM _order.customers WHERE last_ip_address = $1 AND created_at >= NOW() - INTERVAL '5 seconds';

-- name: CompleteOrder :one
UPDATE _order.orders
SET status = 'COMPLETE' 
WHERE id = $1
RETURNING status;

-- name: CreateOrder :one
INSERT INTO _order.orders (
    payment_method,
    paypal_order_id,
    bank_account_id,
    user_id,
    organization_id,
    payment_date,
    job_date,
    job_month,
    job_year,
    distance_round_trip_in_km,
    type,
    status,
    
    total_in_usd
)
VALUES
    (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
    )
RETURNING 
    id,
    bank_account_id,
    user_id,
    organization_id,
    payment_date,
    job_date,
    job_month,
    job_year,
    distance_round_trip_in_km,
    type,
    status,
    
    total_in_usd;

-- name: UpdateOrderStatusToComplete :exec
UPDATE _order.orders SET status = 'COMPLETE' WHERE id = $1;

-- name: UpdateOrderStatusToPending :one
UPDATE _order.orders 
SET payment_date = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000, 
status = 'PENDING' 
WHERE id = $1
RETURNING payment_date, status;

-- name: UpdateOrderStatusToProcessing :exec
UPDATE _order.orders
SET status = 'PROCESSING' 
WHERE id = $1;

-- name: GetOrderTimestamp :one
SELECT 
    id,
    payment_date,
    created_at
FROM _order.orders WHERE id = $1;

-- name: GetOrdersByOrganization :many
SELECT 
    id,
    user_id,
    organization_id,
    payment_date,
    payment_method,
    job_date,
    job_month,
    job_year,
    distance_round_trip_in_km,
    type,
    status,
    total_in_usd,
    paypal_order_id,
    bank_account_id,
    discount_code_id,
    created_at
FROM _order.orders WHERE organization_id = $1;

-- name: GetOrdersByUserId :many
SELECT 
    id,
    organization_id,
    payment_date,
    payment_method,
    job_date,
    job_month,
    job_year,
    distance_round_trip_in_km,
    type,
    status,
    total_in_usd,
    paypal_order_id,
    bank_account_id,
    discount_code_id,
    created_at
FROM _order.orders WHERE user_id = $1
AND payment_date != NULL
ORDER BY payment_date DESC;

-- name: GetOrderByPaypalOrderId :one
SELECT 
    id,
    user_id,
    organization_id,
    payment_date,
    payment_method,
    job_date,
    job_month,
    job_year,
    distance_round_trip_in_km,
    type,
    status,
    
    total_in_usd,
    bank_account_id,
    discount_code_id,
    created_at
FROM _order.orders WHERE paypal_order_id = $1;

-- name: GetOrder :one
SELECT 
    orders.id,
    orders.user_id,
    orders.organization_id,
    orders.payment_date,
    orders.job_date,
    orders.job_month,
    orders.job_year,
    orders.distance_round_trip_in_km,
    orders.type,
    orders.status,
    orders.total_in_usd,
    orders.paypal_order_id,
    orders.bank_account_id,
    orders.created_at

    -- potable_water_orders.quantity_of_water,
    -- potable_water_orders.bottles_required AS quantity_of_potable_water,

    -- irrigation_water_orders.quantity_of_water AS quantity_of_irrigation_water,
    -- irrigation_water_orders.barrels_required,
    -- irrigation_water_orders.quantity_of_barrels,

    -- freight_orders.truck_id,

    -- block_orders.quantity_of_blocks,

    -- block_orders.quantity_of_blocks,

    -- block_orders.quantity_of_blocks,

    -- block_orders.quantity_of_blocks,
FROM _order.orders 
-- LEFT JOIN _order.water_orders ON _order.water_orders.order_id = orders.id 
-- LEFT JOIN _order.freight_orders ON _order.freight_orders.order_id = orders.id 
-- LEFT JOIN _order.block_orders ON _order.block_orders.order_id = orders.id 
WHERE orders.id = $1;

-- name: CreatePotableWaterOrder :one
INSERT INTO _order.potable_water_orders (
    order_id,    
    quantity_of_water,   
    bottles_required
) VALUES (
    $1,
    $2,
    $3
) 
RETURNING 
    id;

-- name: CreateIrrigationWaterOrder :one
INSERT INTO _order.irrigation_water_orders (
    order_id,
    quantity_of_water,
    barrels_required,
    quantity_of_barrels
) VALUES (
    $1,
    $2,
    $3,
    $4
) 
RETURNING 
    id;

-- name: CreateFreightOrder :one
INSERT INTO _order.freight_orders (
    truck_id,
    order_id
) VALUES (
    $1,
    $2
)
RETURNING 
    id;

-- name: CreateRentalVehicleOrder :one
INSERT INTO _order.rental_vehicle_orders (
    order_id,
    rental_vehicle_id,
    start_date,
    end_date
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING 
    id;

-- name: CreateBlockOrder :one
INSERT INTO _order.block_orders (
    order_id,
    quantity_of_blocks
) VALUES (
    $1,
    $2
) 
RETURNING 
    id;

-- name: CreateEggOrder :one
INSERT INTO _order.egg_orders (
    order_id,
    quantity_of_eggs
) VALUES (
    $1,
    $2
) 
RETURNING
    id;

-- name: CreateWaypoint :one
INSERT INTO _order.order_waypoints (
    order_id,
    waypoint,
    address  
) VALUES (
    $1,
    ST_SetSRID(ST_MakePoint(sqlc.arg(longitude)::float, sqlc.arg(latitude)::float), 4326)::geography,
    $2
)
RETURNING 
    id;




-- name: GetFreightOrderByOrderId :one
SELECT 
    id,
    truck_id
FROM _order.freight_orders WHERE order_id = $1;


-- name: GetOrderById :one
SELECT 
    _order.orders.id,
    _order.orders.organization_id,
    _order.orders.user_id,
    _order.orders.paypal_order_id,
    _order.orders.type
FROM _order.orders
WHERE orders.id = $1;

-- name: GetOrdersIncompletedForMonth :many
SELECT 
    id,
    user_id,
    organization_id,
    payment_date,
    payment_method,
    job_date,
    job_month,
    job_year,
    distance_round_trip_in_km,
    type,
    status,
    
    total_in_usd,
    paypal_order_id,
    bank_account_id,
    discount_code_id,
    created_at
FROM _order.orders 
WHERE organization_id = $1 
AND job_month = $2 
AND job_year = $3
AND job_date > EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000 AND status = 'PROCESSING';

-- name: CreateDiscountCode :one
INSERT INTO _order.discount_codes
(
    organization_id,
    code,
    discount_code_type,
    value
)
VALUES
(
    $1,
    $2,
    $3,
    $4
)
RETURNING 
    id;

-- name: GetDiscountCode :one
SELECT 
    id,
    organization_id,
    code,
    discount_code_type,
    value
FROM _order.discount_codes
WHERE code = $1;

-- name: DeleteDiscountCode :exec
DELETE
FROM _order.discount_codes
WHERE id = $1;

-- name: GetAllDiscountCodes :many
SELECT 
    id,
    code,
    discount_code_type,
    value,
    created_at
FROM _order.discount_codes
WHERE id >= sqlc.arg(offset_id)::bigint
AND organization_id = $1
LIMIT $2;

-- name: IsDiscountCodeValid :one
SELECT (
    EXISTS (
        SELECT 1 
        FROM _order.discount_codes
        WHERE _order.discount_codes.code = $1
    )
    AND
    NOT EXISTS (
        SELECT 1 
        FROM _order.orders
        WHERE _order.orders.discount_code_id IN (
            SELECT _order.discount_codes.id from _order.discount_codes WHERE _order.discount_codes.code = $1
        )
    )
) AS is_discount_code_valid;



-- name: GetRentalVehicleByRentTimespan :many
SELECT
    id,
    order_id,
    rental_vehicle_id,
    start_date,
    end_date
FROM _order.rental_vehicle_orders
WHERE end_date <= sqlc.arg(start_date)::bigint OR start_date > sqlc.arg(end_date)::bigint;


-- name: UpdateEnrollment :one
INSERT INTO _order.enrollment (
    vehicle_id,
    freight,
    water_potable,
    water_unpurified,
    ice,
    blocks,
    coconut,
    egg,
    charcoal,
    cooking_gas,
    honey,
    cesspool,
    construction
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13
) 
ON CONFLICT (vehicle_id)
DO UPDATE SET
    freight = $2,
    water_potable = $3,
    water_unpurified = $4,
    ice = $5,
    blocks = $6,
    coconut = $7,
    egg = $8,
    charcoal = $9,
    cooking_gas = $10,
    honey = $11,
    cesspool = $12,
    construction = $13
WHERE _order.enrollment.vehicle_id = $1    
RETURNING 
    id;

-- name: GetEnrollmentByVehicleId :one
SELECT 
    freight,
    water_potable,
    water_unpurified,
    ice,
    blocks,
    coconut,
    egg,
    charcoal,
    cooking_gas,
    honey,
    cesspool,
    construction
FROM _order.enrollment
WHERE vehicle_id = $1;

-- name: SetEnrollmentsDefaultByVehicleId :exec
UPDATE _order.enrollment
SET
    rental = DEFAULT,
    freight = DEFAULT,
    water_potable = DEFAULT,
    water_unpurified = DEFAULT,
    ice = DEFAULT,
    blocks = DEFAULT,
    coconut = DEFAULT,
    egg = DEFAULT,
    charcoal = DEFAULT,
    cooking_gas = DEFAULT,
    honey = DEFAULT,
    cesspool = DEFAULT,
    construction = DEFAULT
WHERE vehicle_id = $1;