-- name: CreateOrganization :one
INSERT INTO organization.organizations
(
    organization_name,
    phone,
    email,
    country,
    currency,
    address,
    duty_free_km,
    max_order_queue_size,
    -- min_selection_time_in_milliseconds,
    address_geog
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6, 
    $7,
    $8,
    -- $9,
    ST_SetSRID(ST_MakePoint(sqlc.arg(address_longitude)::float, sqlc.arg(address_latitude)::float), 4326)::geography
)
RETURNING 
    id;

-- name: UpdateOrganization :exec
UPDATE organization.organizations
SET
    phone = $2,
    email = $3,
    country = $4,
    currency = $5,
    address = $6,
    address_geog = ST_SetSRID(ST_MakePoint(sqlc.arg(address_longitude)::float, sqlc.arg(address_latitude)::float), 4326)::geography,
    duty_free_km = $7,
    max_order_queue_size = $8,
    -- min_selection_time_in_milliseconds = $9,
    is_submitted_for_approval = $9
WHERE id = $1;

-- name: UpdateOrganizationApproval :exec
UPDATE organization.organizations
SET is_approved = TRUE
WHERE id = $1;

-- name: UpdateOrganizationVerifyEmailAndPhone :exec
UPDATE organization.organizations
SET is_email_and_phone_verified = TRUE
WHERE id = $1;

-- name: UpdateOrganizationRating :exec
UPDATE organization.organizations
SET rating = (rating + $2) / 2
WHERE id = $1;

-- name: CreateAccount :one
INSERT INTO organization.accounts
(
    user_id,
    organization_id,
    is_invitation_accepted,
    is_admin,
    is_customer_support,
    is_driver
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING 
    id;

    
-- name: UpdateAccountInvitationAccepted :exec
UPDATE organization.accounts 
SET is_invitation_accepted = TRUE
WHERE id = $1;

-- name: UpdateAccountRoles :exec
UPDATE organization.accounts
SET 
    is_admin = $1,
    is_customer_support = $2,
    is_driver = $3
WHERE id = sqlc.arg(account_id);

-- name: GetOrganization :one
SELECT 
    id,
    avatar_url,
    instagram_handle,
    organization_name,
    phone,
    email,
    country,
    currency,
    address,
    duty_free_km,
    max_order_queue_size,
    ST_X(address_geog::geometry)::float AS address_longitude,
    ST_Y(address_geog::geometry)::float AS address_latitude,
    is_email_and_phone_verified,
    is_submitted_for_approval,
    rating,
    is_approved
FROM organization.organizations WHERE id = $1 AND is_deleted = FALSE;

-- name: CreateInstagramHandleVerification :one
INSERT INTO organization.instagram_handle_verification
(
    organization_id
)
VALUES
(
    $1
)
ON CONFLICT (organization_id)
DO UPDATE SET
    instagram_handle_verification_code = DEFAULT,
    instagram_handle_verification_code_timestamp = DEFAULT,
    failed_attempts = DEFAULT
RETURNING 
    id,
    instagram_handle_verification_code,
    instagram_handle_verification_code_timestamp,
    failed_attempts;

-- name: IncrementInstagramHandleVerificationFailedAttempt :exec
UPDATE organization.instagram_handle_verification
SET failed_attempts = failed_attempts + 1
WHERE id = $1;

-- name: GetInstagramHandleVerificationByOrganizationId :one
SELECT 
    id,
    instagram_handle_verification_code,
    instagram_handle_verification_code_timestamp,
    failed_attempts,
    created_at
FROM organization.instagram_handle_verification
WHERE organization_id = $1;

-- name: UpdateInstagramHandleAndAvatarUrl :exec
UPDATE organization.organizations
SET 
    instagram_handle = $1,
    avatar_url = $2
WHERE id = sqlc.arg(organization_id);

-- name: GetOrganizationApprovals :many
SELECT 
    id,
    organization_id,
    system_administrator_id,
    created_at
FROM organization.organization_approvals WHERE organization_id = $1 AND is_deleted = FALSE;

-- name: GetAccount :one
SELECT 
    organization.accounts.user_id,
    organization.accounts.organization_id,
    organization.accounts.is_invitation_accepted,
    organization.accounts.is_admin,
    organization.accounts.is_customer_support,
    organization.accounts.is_driver,
    _user.users.first_name,
    _user.users.last_name,
    _user.users.phone,
    _user.users.email,
    _user.users.created_at
FROM organization.accounts
JOIN _user.users ON organization.accounts.user_id = _user.users.id
WHERE organization.accounts.id = $1;

-- name: GetAccountByUserIdAndOrganizationId :one
SELECT 
    organization.accounts.user_id,
    organization.accounts.organization_id,
    organization.accounts.is_invitation_accepted,
    organization.accounts.is_admin,
    organization.accounts.is_customer_support,
    organization.accounts.is_driver,
    _user.users.first_name,
    _user.users.last_name,
    _user.users.phone,
    _user.users.email,
    _user.users.created_at
FROM organization.accounts
JOIN _user.users ON organization.accounts.user_id = _user.users.id
WHERE organization.accounts.user_id = $1 AND organization.accounts.organization_id = $2;


-- name: GetOrganizationsPendingByUserId :many
SELECT 
    organization.organizations.id,
    organization.organizations.organization_name,
    organization.organizations.phone,
    organization.organizations.email,
    organization.organizations.country,
    organization.organizations.currency,
    organization.organizations.address,
    organization.organizations.duty_free_km,
    organization.organizations.max_order_queue_size,
    ST_X(organization.organizations.address_geog::geometry)::float AS address_longitude,
    ST_Y(organization.organizations.address_geog::geometry)::float AS address_latitude,
    organization.accounts.id as accountId,
    organization.accounts.user_id,
    organization.accounts.organization_id,
    organization.accounts.is_invitation_accepted,
    organization.accounts.is_admin,
    organization.accounts.is_customer_support,
    organization.accounts.is_driver

FROM organization.organizations 
JOIN organization.accounts
ON organization.organizations.id = organization.accounts.organization_id
WHERE organization.accounts.user_id = sqlc.arg(user_id)
AND organization.organizations.is_approved = FALSE;

-- -- name: GetSalesUnits :many
-- SELECT
--     id,
--     quantity
-- FROM organization.sales_units
-- WHERE organization_id = $1 AND order_type = $2;

-- -- name: GetAllSalesUnits :many
-- SELECT
--     id,
--     quantity
-- FROM organization.sales_units
-- WHERE organization_id = $1;

-- name: GetFreightConfig :one
SELECT 
    id,
    organization_id,
    freight_price_min,
    ST_X(solid_waste_dump_site_geog::geometry)::float AS solid_waste_dump_site_longitude,
    ST_Y(solid_waste_dump_site_geog::geometry)::float AS soplid_waste_dump_site_latitude,
    ST_X(liquid_waste_dump_site_geog::geometry)::float AS liquid_waste_dump_site_longitude,
    ST_Y(liquid_waste_dump_site_geog::geometry)::float AS liquid_waste_dump_site_latitude,
    is_enabled
FROM organization.freight_config
WHERE organization_id = $1;

-- name: GetPotableWaterSalesConfig :one
SELECT 
    id,
    organization_id,
    preferred_unit,
    water_price_min,
    water_price_per_km,
    water_price,
    bottle_cost,
    available_supply,    
    min_order_quantity,
    is_enabled
FROM organization.potable_water_sales_config
WHERE organization_id = $1;

-- name: GetIrrigationWaterSalesConfig :one
SELECT 
    id,
    organization_id,
    preferred_unit,
    water_price_min,
    water_price_per_km,
    water_price,
    barrel_cost,
    available_supply,    
    min_order_quantity,
    is_enabled
FROM organization.irrigation_water_sales_config
WHERE organization_id = $1;

-- name: GetBlockSalesConfig :one
SELECT 
    id,
    organization_id,
    preferred_unit,
    block_price_min,
    block_price,
    block_price_per_km,
    available_supply,    
    min_order_quantity,
    is_enabled
FROM organization.block_sales_config
WHERE organization_id = $1;

-- name: GetHoneySalesConfig :one
SELECT 
    id,
    organization_id,
    preferred_unit,
    honey_price_min,
    honey_price,
    honey_price_per_km,
    available_supply,    
    min_order_quantity,
    is_enabled
FROM organization.honey_sales_config
WHERE organization_id = $1;

-- name: GetEggSalesConfig :one
SELECT 
    id,
    organization_id,
    preferred_unit,
    egg_price_min,
    egg_price,
    egg_price_per_km,
    available_supply,    
    min_order_quantity,
    is_enabled
FROM organization.egg_sales_config
WHERE organization_id = $1;

-- name: GetCharcoalSalesConfig :one
SELECT 
    id,
    organization_id,
    preferred_unit,
    charcoal_price_min,
    charcoal_price,
    charcoal_price_per_km,
    available_supply,    
    min_order_quantity,
    is_enabled
FROM organization.charcoal_sales_config
WHERE organization_id = $1;

-- name: GetCoconutSalesConfig :one
SELECT 
    id,
    organization_id,
    preferred_unit,
    coconut_price_min,
    coconut_price,
    coconut_price_per_km,
    available_supply,    
    min_order_quantity,
    is_enabled
FROM organization.coconut_sales_config
WHERE organization_id = $1;

-- name: GetIceSalesConfig :one
SELECT 
    id,
    organization_id,
    preferred_unit,
    ice_price_min,
    ice_price,
    ice_price_per_km,
    available_supply,    
    min_order_quantity,
    is_enabled
FROM organization.ice_sales_config
WHERE organization_id = $1;


-- name: GetOrganizationMaxOrderQueueLength :one
SELECT max_order_queue_size FROM organization.organizations WHERE id = $1;

-- name: UpsertPayPalCredential :one
INSERT INTO organization.paypal_credentials (
    organization_id,
    email,
    subscription_id,
    rest_sdk_client_id_encrypted,
    rest_sdk_secret_encrypted
) VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (organization_id) 
DO UPDATE SET
    email = EXCLUDED.email,
    subscription_id = EXCLUDED.subscription_id,
    rest_sdk_client_id_encrypted = EXCLUDED.rest_sdk_client_id_encrypted,
    rest_sdk_secret_encrypted = EXCLUDED.rest_sdk_secret_encrypted
RETURNING id, organization_id, email, subscription_id, rest_sdk_client_id_encrypted, rest_sdk_secret_encrypted;

-- name: GetPayPalCredentialByOrganizationID :one
SELECT 
    id, 
    organization_id, 
    email, 
    subscription_id, 
    rest_sdk_client_id_encrypted, 
    rest_sdk_secret_encrypted
FROM organization.paypal_credentials
WHERE organization_id = $1;

-- name: DeletePayPalCredentialByOrganizationID :exec
DELETE FROM organization.paypal_credentials
WHERE organization_id = $1;

-- name: CreateBitcoinCredential :one
INSERT INTO organization.bitcoin_credentials (
    organization_id
) VALUES (
    $1
)
RETURNING id, organization_id;

-- name: DeleteBitcoinCredentialByOrganizationID :exec
DELETE FROM organization.bitcoin_credentials
WHERE organization_id = $1;

-- name: CreateBitcoinLightningCredential :one
INSERT INTO organization.bitcoin_lightning_credentials (
    organization_id, 
    lightning_address_encrypted
) VALUES (
    $1, $2
)
RETURNING id, organization_id, lightning_address_encrypted;

-- name: DeleteBitcoinLightningCredentialByOrganizationID :exec
DELETE FROM organization.bitcoin_lightning_credentials
WHERE organization_id = $1;

-- name: GetBitcoinCredentialByOrganizationID :one
SELECT 
    id, 
    organization_id
FROM organization.bitcoin_credentials
WHERE organization_id = $1;

-- name: GetBitcoinLightningCredentialByOrganizationID :one
SELECT 
    id, 
    organization_id, 
    lightning_address_encrypted
FROM organization.bitcoin_lightning_credentials
WHERE organization_id = $1;

-- name: GetAccounts :many
SELECT
_user.users.id as user_id,
_user.users.first_name,
_user.users.last_name,
_user.users.phone,
_user.users.email,
_user.users.password_hash,
_user.users.google_user_id,
_user.users.oauth_refresh_token_encrypted,
_user.users.last_ip_address,
_user.users.created_at,
organization.accounts.id,
organization.accounts.user_id,
organization.accounts.organization_id,
organization.accounts.is_admin,
organization.accounts.is_customer_support,
organization.accounts.is_driver
FROM organization.accounts 
JOIN _user.users ON organization.accounts.user_id = _user.users.id
WHERE organization_id = $1;

-- name: GetDrivers :many
SELECT
_user.users.id as user_id,
_user.users.first_name,
_user.users.last_name,
_user.users.phone,
_user.users.email,
_user.users.password_hash,
_user.users.google_user_id,
_user.users.oauth_refresh_token_encrypted,
_user.users.last_ip_address,
_user.users.created_at,
organization.accounts.id,
organization.accounts.user_id,
organization.accounts.organization_id,
organization.accounts.is_admin,
organization.accounts.is_customer_support,
organization.accounts.is_driver
FROM organization.accounts 
JOIN _user.users ON organization.accounts.user_id = _user.users.id
WHERE organization_id = $1
AND is_driver = TRUE;

-- name: GetAccountsCount :one
SELECT
    COUNT(*)
FROM organization.accounts
WHERE organization_id = $1;

-- name: GetAccountsConfirmed :many
SELECT
_user.users.id,
_user.users.first_name,
_user.users.last_name,
_user.users.phone,
_user.users.email,
_user.users.password_hash,
_user.users.google_user_id,
_user.users.oauth_refresh_token_encrypted,
_user.users.last_ip_address,
_user.users.created_at,
organization.accounts.id,
organization.accounts.user_id,
organization.accounts.organization_id,
organization.accounts.is_admin,
organization.accounts.is_customer_support,
organization.accounts.is_driver
FROM organization.accounts 
JOIN _user.users ON organization.accounts.user_id = _user.users.id
WHERE organization_id = $1
AND is_invitation_accepted = TRUE;

-- name: GetAllBankAccountInfo :many
SELECT
    id,
    organization_id,
    bank,
    bank_branch,
    account_name,
    type,
    ownership,
    account_number,
    swift_code,
    routing_number,
    recipient_address
FROM organization.bank_accounts WHERE organization_id = $1;

-- name: GetBankAccountInfo :one
SELECT
id,
organization_id,
bank,
bank_branch,
type,
ownership,
account_number,
account_name,
swift_code
FROM organization.bank_accounts WHERE id = $1;

-- name: GetOrganizationByTruckId :one
SELECT
    organization.organizations.id,
    organization.organizations.organization_name,
    organization.organizations.phone,
    organization.organizations.email,
    organization.organizations.country,
    organization.organizations.currency,
    organization.organizations.duty_free_km,
    organization.organizations.address,
    organization.organizations.rating,
    organization.organizations.max_order_queue_size,
    organization.organizations.is_approved,
    ST_X(organization.organizations.address_geog::geometry)::float AS address_longitude,
    ST_Y(organization.organizations.address_geog::geometry)::float AS address_latitude
FROM organization.organizations 
JOIN logistic.trucks ON organization.organizations.id = logistic.trucks.organization_id
WHERE logistic.trucks.id = sqlc.arg(truck_id)::bigint AND organization.organizations.is_deleted = FALSE;


-- name: CreateOrganizationApproval :one
INSERT INTO organization.organization_approvals
(
    organization_id,
    system_administrator_id
)
VALUES
(
    $1,
    $2
)
RETURNING 
    id;

-- name: GetOrganizationsForUser :many
SELECT 
    organization.organizations.id,
    organization.organizations.avatar_url,
    organization.organizations.organization_name,
    organization.organizations.rating,
    organization.organizations.is_approved
FROM organization.organizations
JOIN organization.accounts
ON organization.organizations.id = organization.accounts.organization_id
WHERE organization.organizations.is_deleted = FALSE
AND organization.organizations.id >= sqlc.arg(offset_id)::bigint
AND organization.organizations.organization_name ILIKE '%' || sqlc.arg(filter_text)::text || '%'
AND organization.accounts.user_id = sqlc.arg(user_id)::bigint
LIMIT sqlc.arg(limitation)::bigint;


-- name: CreateOrganizationVerificationCode :one
INSERT INTO organization.organization_verification_code
    (
        organization_id,
        phone_verification_codes_id,
        email_verification_codes_id
    )
VALUES
    (
        $1,
        $2,
        $3
    )
ON CONFLICT (organization_id)
DO UPDATE SET 
        phone_verification_codes_id = $2,
        email_verification_codes_id = $3
RETURNING
    id;

-- name: GetOrganizationVerificationCodeByOrganizationId :one
SELECT
    id,
    phone_verification_codes_id,
    email_verification_codes_id,
    created_at
FROM organization.organization_verification_code
WHERE organization_id = $1;


-- name: CreateOrganizationMailList :one
INSERT INTO organization.organization_mail_list
    (
        organization_id,
        email
    )
VALUES
    (
        $1,
        $2
    )
RETURNING
    id;


-- name: GetOrganizationMailListByOrganizationId :many
SELECT
    id,
    email
FROM organization.organization_mail_list
WHERE organization_id = $1;

-- name: DeleteOrganizationById :exec
DELETE FROM organization.organizations
WHERE id = $1;

-- name: DeleteOrganizations :exec
DELETE FROM organization.organizations;

-- name: DeleteOrganizationMailList :exec
DELETE FROM organization.organization_mail_list
WHERE id = $1;



-- name: UpdateOrganizationFreightConfig :exec
INSERT INTO organization.freight_config
(
    organization_id,
    freight_price_min,
    solid_waste_dump_site_geog,
    liquid_waste_dump_site_geog,
    is_enabled
)
VALUES
(
    $1,
    $2,
    ST_SetSRID(ST_MakePoint(sqlc.arg(solid_waste_dump_site_longitude)::float, sqlc.arg(solid_waste_dump_site_latitude)::float), 4326)::geography,
    ST_SetSRID(ST_MakePoint(sqlc.arg(liquid_waste_dump_site_longitude)::float, sqlc.arg(liquid_waste_dump_site_latitude)::float), 4326)::geography,
    $3
)
ON CONFLICT (organization_id)
DO UPDATE SET
    freight_price_min = $2,
    solid_waste_dump_site_geog = ST_SetSRID(ST_MakePoint(sqlc.arg(solid_waste_dump_site_longitude)::float, sqlc.arg(solid_waste_dump_site_latitude)::float), 4326)::geography,
    liquid_waste_dump_site_geog = ST_SetSRID(ST_MakePoint(sqlc.arg(liquid_waste_dump_site_longitude)::float, sqlc.arg(liquid_waste_dump_site_latitude)::float), 4326)::geography,
    is_enabled = $3
WHERE organization_id = $1
RETURNING 
    id;

-- name: UpdateOrganizationPotableWaterSalesConfig :exec
INSERT INTO organization.potable_water_sales_config
(
    organization_id,
    preferred_unit,
    water_price_min,
    water_price_per_km,
    water_price,
    bottle_cost,
    available_supply,
    min_order_quantity,
    is_enabled
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
ON CONFLICT (organization_id)
DO UPDATE SET
    preferred_unit = $2,
    water_price_min = $3,
    water_price_per_km = $4,
    water_price_per_gallon = $5,
    bottle_cost = $6,
    available_supply = $7,
    min_order_quantity = $8,
    is_enabled = $9
WHERE organization_id = $1 
RETURNING 
    id;

-- name: UpdateOrganizationIrrigationWaterSalesConfig :exec
INSERT INTO organization.irrigation_water_sales_config
(
    organization_id,
    preferred_unit,
    water_price_min,
    water_price_per_km,
    water_price,
    barrel_cost,
    available_supply,
    min_order_quantity,
    is_enabled
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
ON CONFLICT (organization_id)
DO UPDATE SET
    preferred_unit = $2,
    water_price_min = $3,
    water_price_per_km = $4,
    bottle_cost = $5,
    barrel_cost = $6,
    available_supply = $7,
    min_order_quantity = $8,
    is_enabled = $9
WHERE organization_id = $1 
RETURNING 
    id;

-- name: UpdateOrganizationBlocksSalesConfig :exec
INSERT INTO organization.block_sales_config
(
    organization_id,
    block_price_min,
    block_price,
    block_price_per_km,
    available_supply,
    min_order_quantity,
    is_enabled
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
ON CONFLICT (organization_id)
DO UPDATE SET
    block_price_min = $2,
    block_price = $3,
    block_price_per_km = $4,
    available_supply = $5,
    min_order_quantity = $6,
    is_enabled = $7
WHERE organization_id = $1 
RETURNING 
    id;

-- name: UpdateOrganizationIceSalesConfig :exec
INSERT INTO organization.ice_sales_config
(
    organization_id,
    ice_price_min,
    ice_price,
    ice_price_per_km,
    available_supply,
    min_order_quantity,
    is_enabled
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
ON CONFLICT (organization_id)
DO UPDATE SET
    ice_price_min = $2,
    ice_price = $3,
    ice_price_per_km = $4,
    available_supply = $5,
    min_order_quantity = $6,
    is_enabled = $7
WHERE organization_id = $1 
RETURNING 
    id;

-- name: UpdateOrganizationHoneySalesConfig :exec
INSERT INTO organization.honey_sales_config
(
    organization_id,
    honey_price_min,
    honey_price,
    honey_price_per_km,
    available_supply,
    min_order_quantity,
    is_enabled
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
ON CONFLICT (organization_id)
DO UPDATE SET
    honey_price_min = $2,
    honey_price = $3,
    honey_price_per_km = $4,
    available_supply = $5,
    min_order_quantity = $6,
    is_enabled = $7
WHERE organization_id = $1 
RETURNING 
    id;

-- name: UpdateOrganizationEggSalesConfig :exec
INSERT INTO organization.egg_sales_config
(
    organization_id,
    egg_price_min,
    egg_price,
    egg_price_per_km,
    available_supply,
    min_order_quantity,
    is_enabled
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
ON CONFLICT (organization_id)
DO UPDATE SET
    egg_price_min = $2,
    egg_price = $3,
    egg_price_per_km = $4,
    available_supply = $5,
    min_order_quantity = $6,
    is_enabled = $7
WHERE organization_id = $1 
RETURNING 
    id;

-- name: UpdateOrganizationCoconutSalesConfig :exec
INSERT INTO organization.coconut_sales_config
(
    organization_id,
    coconut_price_min,
    coconut_price,
    coconut_price_per_km,
    available_supply,
    min_order_quantity,
    is_enabled
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
ON CONFLICT (organization_id)
DO UPDATE SET
    coconut_price_min = $2,
    coconut_price = $3,
    coconut_price_per_km = $4,
    available_supply = $5,
    min_order_quantity = $6,
    is_enabled = $7
WHERE organization_id = $1 
RETURNING 
    id;


-- name: GetOrganizationSupportedOrderTypes :many
WITH target_orgs AS (
    -- Get all relevant orgs and their public vehicle status once
    SELECT 
        o.id AS org_id,
        t.is_public
    FROM organization.organizations o
    JOIN logistic.trucks t ON o.id = t.organization_id
    JOIN logistic.vehicles v ON t.vehicle_id = v.id
    WHERE o.country = $1
),
rental_orgs AS (
    -- Rental has a different join path, so we isolate it
    SELECT o.id AS org_id
    FROM organization.organizations o
    JOIN logistic.rental_vehicles rv ON o.id = rv.organization_id
    WHERE o.country = $1
)
SELECT ARRAY(

    -- TODO: Add additional order types

    -- Check for FREIGHT
    SELECT 'FREIGHT'::shared.order_type 
    WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.freight_config c ON target_orgs.org_id = c.organization_id WHERE organization.freight_config.is_enabled = TRUE AND target_orgs.is_public = TRUE)
    UNION ALL
    -- Check for VEHICLE_RENTAL
    SELECT 'VEHICLE_RENTAL'::shared.order_type 
    WHERE EXISTS (SELECT 1 FROM rental_orgs JOIN logistic.vehicles v ON TRUE -- logic depends on your schema here
                  WHERE EXISTS (SELECT 1 FROM logistic.rental_vehicles rv WHERE rv.organization_id = rental_orgs.org_id))
    UNION ALL
    -- Check for WATER_PURIFIED
    SELECT 'WATER_PURIFIED'::shared.order_type 
    WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.potable_water_sales_config c ON target_orgs.org_id = c.organization_id WHERE organization.potable_water_sales_config.is_enabled = TRUE AND  target_orgs.is_public = TRUE AND organization.potable_water_sales_config.water_price_per_gallon != null)
    UNION ALL
    -- Check for WATER_UNPURIFIED
    SELECT 'WATER_UNPURIFIED'::shared.order_type 
    WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.irrigation_water_sales_config c ON target_orgs.org_id = c.organization_id WHERE organization.irrigation_water_sales_config.is_enabled = TRUE AND target_orgs.is_public = TRUE AND organization.irrigation_water_sales_config.water_price_per_gallon != null)
    UNION ALL
    -- Check for BLOCK
    SELECT 'BLOCK'::shared.order_type 
    WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.block_sales_config c ON target_orgs.org_id = c.organization_id WHERE organization.block_sales_config.is_enabled = TRUE AND target_orgs.is_public = TRUE)
    UNION ALL
    -- Check for COCONUT
    SELECT 'COCONUT'::shared.order_type 
    WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.coconut_sales_config c ON target_orgs.org_id = c.organization_id WHERE organization.coconut_sales_config.is_enabled = TRUE AND target_orgs.is_public = TRUE)
    UNION ALL
    -- Check for EGG
    SELECT 'EGG'::shared.order_type 
    WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.egg_sales_config c ON target_orgs.org_id = c.organization_id WHERE organization.egg_sales_config.is_enabled = TRUE AND target_orgs.is_public = TRUE)
)::shared.order_type AS supported_country_order_types;


-- -- name: GetOrganizationUnsupportedOrderTypes :many
-- WITH target_orgs AS (
--     -- Step 1: Identify all orgs in the country once
--     SELECT 
--         o.id AS org_id,
--         v.is_public
--     FROM organization.organizations o
--     JOIN logistic.trucks t ON o.id = t.organization_id
--     JOIN logistic.vehicles v ON t.vehicle_id = v.id
--     WHERE o.country = $1
-- ),
-- supported_types AS (
--     -- Step 2: Build the list of what IS supported (using the efficient refactor)
--     SELECT 'FREIGHT'::shared.order_type WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.freight_config c ON target_orgs.org_id = c.organization_id WHERE target_orgs.is_public = TRUE)
--     UNION ALL
--     SELECT 'VEHICLE_RENTAL'::shared.order_type WHERE EXISTS (SELECT 1 FROM organization.organizations o JOIN logistic.rental_vehicles rv ON o.id = rv.organization_id WHERE o.country = $1)
--     UNION ALL
--     SELECT 'WATER'::shared.order_type WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.water_sales_config c ON target_orgs.org_id = c.organization_id WHERE target_orgs.is_public = TRUE)
--     UNION ALL
--     SELECT 'BLOCK'::shared.order_type WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.block_sales_config c ON target_orgs.org_id = c.organization_id WHERE target_orgs.is_public = TRUE)
--     UNION ALL
--     SELECT 'COCONUT'::shared.order_type WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.coconut_sales_config c ON target_orgs.org_id = c.organization_id WHERE target_orgs.is_public = TRUE)
--     UNION ALL
--     SELECT 'EGG'::shared.order_type WHERE EXISTS (SELECT 1 FROM target_orgs JOIN organization.egg_sales_config c ON target_orgs.org_id = c.organization_id WHERE target_orgs.is_public = TRUE)
-- )
-- SELECT ARRAY(
--     -- Step 3: All possible enum values minus the supported ones
--     SELECT unnest(enum_range(NULL::shared.order_type))
--     EXCEPT
--     SELECT * FROM supported_types
-- ) AS unsupported_country_order_types;
