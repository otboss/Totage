package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) CreateTruckHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateTruckCommand],
) (*connect.Response[endpoints.CreateTruckResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	vehicle, err := db_logistic.Query.GetVehicle(ctx, req.Msg.VehicleId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	organization, err := db_organization.Query.GetOrganization(ctx, vehicle.OrganizationID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: organization.ID,
	})
	if err != nil {
		return nil, errors.New("account not found")
	}

	if account.IsAdmin == false {
		return nil, errors.New("insufficient account permission")
	}

	// exchangeRateUsd, err := clients.CurrencyApi.GetExchangeRateUsd(string(organization.Currency))
	// if err != nil {
	// 	return nil, connect.NewError(
	// 		connect.CodeInternal,
	// 		errors.New("internal server error"),
	// 	)
	// }

	var bedHeight pgtype.Float8
	if req.Msg.Truck.BedSize.HeightMm != nil {
		bedHeight = pgtype.Float8{
			Float64: float64(*req.Msg.Truck.BedSize.HeightMm),
			Valid:   true,
		}
	}

	truckProperty, truckPropertyExists := common.FreightTruckProperty_name[int32(req.Msg.Truck.Property)]
	if truckPropertyExists == false {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid truck property provided"),
		)
	}

	// var truckProperty common.FreightTruckProperty = common.FreightTruckProperty(common.FreightTruckProperty_value[req.Msg.Truck.Property])

	truckPropertyFromDb, err := db_logistic.Query.GetTruckPropertyByProperty(ctx, truckProperty)
	if err == pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid truck property provided"),
		)
	}
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	// truckProperty, err := local_utils.TruckPropertyMap.ProtoTruckPropertyToDbTruckProperty(req.Msg.Truck.Property)
	// if err != nil {
	// 	return nil, connect.NewError(
	// 		connect.CodeInternal,
	// 		errors.New("internal server error"),
	// 	)
	// }

	truckId, err := db_logistic.Query.CreateTruck(ctx, db_logistic.CreateTruckParams{
		VehicleID:         vehicle.ID,
		Cost:              float64(req.Msg.Truck.Cost),
		FreightPricePerKm: float64(req.Msg.Truck.FreightPricePerKm),
		Property:          truckPropertyFromDb.ID,
		BedLengthMm:       float64(req.Msg.Truck.BedSize.LengthMm),
		BedWidthMm:        float64(req.Msg.Truck.BedSize.WidthMm),
		BedHeightMm:       bedHeight,
		IsPublic:          req.Msg.Truck.IsPublic,
		Longitude:         organization.AddressLongitude,
		Latitude:          organization.AddressLatitude,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.CreateTruckResponse{
		TruckId: truckId,
	}), nil
}
