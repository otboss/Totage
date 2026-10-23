package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"

	"connectrpc.com/connect"
)

func (s *Handlers) GetTruckHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetTruckQuery],
) (*connect.Response[endpoints.GetTruckResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	truck, err := db_logistic.Query.GetTruckByVehicleId(ctx, req.Msg.VehicleId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	vehicle, err := db_logistic.Query.GetVehicle(ctx, truck.VehicleID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: vehicle.OrganizationID,
	})
	if err != nil {
		return nil, errors.New("account not found")
	}

	if account.IsAdmin == false {
		return nil, errors.New("user not permitted")
	}

	var truckBedHeight *float32 = nil
	if truck.BedHeightMm.Valid {
		bedHeight := float32(truck.BedHeightMm.Float64)
		truckBedHeight = &bedHeight
	}

	truckBedSize := common.BedSize{
		LengthMm: float32(truck.BedLengthMm),
		WidthMm:  float32(truck.BedWidthMm),
		HeightMm: truckBedHeight,
	}

	truckPropertyFromDb, err := db_logistic.Query.GetTruckPropertyByPropertyId(ctx, truck.Property)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	truckPropertyNum, exists := common.FreightTruckProperty_value[truckPropertyFromDb.Property]
	if exists == false {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.GetTruckResponse{
		Truck: &common.Truck{
			Id:                truck.ID,
			VehicleId:         truck.VehicleID,
			BedSize:           &truckBedSize,
			Property:          common.FreightTruckProperty(truckPropertyNum),
			FreightPricePerKm: float32(truck.FreightPricePerKm),
			Cost:              float32(truck.Cost),
			IsPublic:          truck.IsPublic,
		},
	}), nil
}
