#!/bin/bash
cat \
"shared/domain/entities/shared.schema.sql" \
"modules/user/domain/entities/user.schema.sql" \
"modules/payment/domain/entities/payment.schema.sql" \
"modules/organization/domain/entities/organization.schema.sql" \
"modules/logistic/domain/entities/logistic.schema.sql" \
"modules/order/domain/entities/order.schema.sql"