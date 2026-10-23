package application

import (
	"context"
	"errors"
	"fmt"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) CreateOrganizationHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateOrganizationCommand],
) (*connect.Response[endpoints.CreateOrganizationResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	_, err = db_user.Query.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var organization *db_organization.GetOrganizationRow
	if req.Msg.OrganizationId != nil {
		account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
			UserID:         claims.UserID,
			OrganizationID: *req.Msg.OrganizationId,
		})
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		if account.IsAdmin == false {
			return nil, connect.NewError(
				connect.CodePermissionDenied,
				errors.New("action not permitted"),
			)
		}

		org, err := db_organization.Query.GetOrganization(ctx, *req.Msg.OrganizationId)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
		organization = &org
	}

	if organization != nil && organization.IsApproved && req.Msg.OrganizationName != organization.OrganizationName {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("cannot update name of submitted/approved organization"),
		)
	}

	var maxOrderQueueSize pgtype.Int2 = pgtype.Int2{
		Int16: 0,
		Valid: false,
	}
	if req.Msg.MaxOrderQueueSize != nil {
		maxOrderQueueSize = pgtype.Int2{
			Int16: int16(*req.Msg.MaxOrderQueueSize),
			Valid: true,
		}
	}

	geocodingResult, err := clients.GoogleCloudApi.GeocodeAddress(fmt.Sprintf("%s, %s", req.Msg.Address, req.Msg.Country.String()))
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}
	if len(geocodingResult.Results) == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid address provided"),
		)
	}

	if organization == nil {
		organizationId, err := db_organization.Query.CreateOrganization(ctx, db_organization.CreateOrganizationParams{
			OrganizationName:  req.Msg.OrganizationName,
			Phone:             req.Msg.Phone,
			Email:             req.Msg.Email,
			Country:           db_organization.SharedCountryCode(req.Msg.Country.String()),
			Currency:          db_organization.SharedCurrency(req.Msg.Currency.String()),
			MaxOrderQueueSize: maxOrderQueueSize,
			Address:           req.Msg.Address,
			AddressLongitude:  geocodingResult.Results[0].Geometry.Location.Lng,
			AddressLatitude:   geocodingResult.Results[0].Geometry.Location.Lat,
			DutyFreeKm:        req.Msg.DutyFreeKm,
		})
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		_, err = db_organization.Query.CreateAccount(ctx, db_organization.CreateAccountParams{
			UserID:               claims.UserID,
			OrganizationID:       organizationId,
			IsInvitationAccepted: true,
			IsAdmin:              true,
			IsCustomerSupport:    true,
			IsDriver:             true,
		})
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		org, err := db_organization.Query.GetOrganization(ctx, organizationId)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
		organization = &org
	} else {
		err = db_organization.Query.UpdateOrganization(ctx, db_organization.UpdateOrganizationParams{
			ID:                     organization.ID,
			Phone:                  req.Msg.Phone,
			Email:                  req.Msg.Email,
			Country:                db_organization.SharedCountryCode(req.Msg.Country.String()),
			Currency:               db_organization.SharedCurrency(req.Msg.Country.String()),
			MaxOrderQueueSize:      maxOrderQueueSize,
			Address:                req.Msg.Address,
			AddressLongitude:       geocodingResult.Results[0].Geometry.Location.Lng,
			AddressLatitude:        geocodingResult.Results[0].Geometry.Location.Lat,
			DutyFreeKm:             req.Msg.DutyFreeKm,
			IsSubmittedForApproval: true,
		})
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
	}

	if organization.IsApproved || organization.IsSubmittedForApproval {
		if organization.OrganizationName != req.Msg.OrganizationName {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("cannot update name of submitted/approved organization"),
			)
		}
	}

	//
	//
	//

	configSet := false

	// if req.Msg.FreightConfig != nil {
	// 	err = db_organization.Query.UpdateOrganizationFreightConfig(ctx, db_organization.UpdateOrganizationFreightConfigParams{
	// 		OrganizationID:  organization.ID,
	// 		FreightPriceMin: req.Msg.FreightConfig.MinimumJobCost,
	// 		IsEnabled: pgtype.Bool{
	// 			Bool:  req.Msg.FreightConfig.IsEnabled,
	// 			Valid: true,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, connect.NewError(
	// 			connect.CodeInternal,
	// 			errors.New("internal server error"),
	// 		)
	// 	}
	// 	configSet = true
	// }

	// if req.Msg.PotableWaterDeliveryConfig != nil {
	// 	// TODO:
	// 	err = db_organization.Query.UpdateOrganizationPotableWaterSalesConfig(ctx, db_organization.UpdateOrganizationPotableWaterSalesConfigParams{
	// 		OrganizationID:  organization.ID,
	// 		WaterPriceMin:   req.Msg.PotableWaterDeliveryConfig.MinimumJobCost,
	// 		WaterPricePerKm: req.Msg.PotableWaterDeliveryConfig.DutyCostPerKm,
	// 		WaterPricePerGallon: pgtype.Float8{
	// 			Float64: req.Msg.PotableWaterDeliveryConfig.WaterCost,
	// 			Valid:   true,
	// 		},
	// 		BottleCost: pgtype.Float8{
	// 			Float64: req.Msg.PotableWaterDeliveryConfig.BottleCost,
	// 			Valid:   true,
	// 		},
	// 		IsEnabled: pgtype.Bool{
	// 			Bool:  req.Msg.PotableWaterDeliveryConfig.IsEnabled,
	// 			Valid: true,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, connect.NewError(
	// 			connect.CodeInternal,
	// 			errors.New("internal server error"),
	// 		)
	// 	}
	// 	configSet = true
	// }

	// if req.Msg.IrrigationWaterDeliveryConfig != nil {
	// 	// TODO:
	// 	err = db_organization.Query.UpdateOrganizationIrrigationWaterSalesConfig(ctx, db_organization.UpdateOrganizationIrrigationWaterSalesConfigParams{
	// 		OrganizationID:  organization.ID,
	// 		WaterPriceMin:   req.Msg.IrrigationWaterDeliveryConfig.MinimumJobCost,
	// 		WaterPricePerKm: req.Msg.IrrigationWaterDeliveryConfig.DutyCostPerKm,
	// 		WaterPricePerGallon: pgtype.Float8{
	// 			Float64: req.Msg.IrrigationWaterDeliveryConfig.WaterCost,
	// 			Valid:   true,
	// 		},
	// 		BarrelCost: pgtype.Float8{
	// 			Float64: req.Msg.IrrigationWaterDeliveryConfig.BarrelCost,
	// 			Valid:   true,
	// 		},
	// 		IsEnabled: pgtype.Bool{
	// 			Bool:  req.Msg.IrrigationWaterDeliveryConfig.IsEnabled,
	// 			Valid: true,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, connect.NewError(
	// 			connect.CodeInternal,
	// 			errors.New("internal server error"),
	// 		)
	// 	}
	// 	configSet = true
	// }

	// if req.Msg.BlocksDeliveryConfig != nil {
	// 	err = db_organization.Query.UpdateOrganizationBlocksSalesConfig(ctx, db_organization.UpdateOrganizationBlocksSalesConfigParams{
	// 		OrganizationID: organization.ID,
	// 		// PreferredUnit: nil,
	// 		BlockPriceMin:   req.Msg.BlocksDeliveryConfig.MinimumJobCost,
	// 		BlockPrice:      req.Msg.BlocksDeliveryConfig.BlockCost,
	// 		BlockPricePerKm: req.Msg.BlocksDeliveryConfig.DutyCostPerKm,
	// 		IsEnabled: pgtype.Bool{
	// 			Bool:  req.Msg.BlocksDeliveryConfig.IsEnabled,
	// 			Valid: true,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, connect.NewError(
	// 			connect.CodeInternal,
	// 			errors.New("internal server error"),
	// 		)
	// 	}
	// 	configSet = true
	// }

	// if req.Msg.EggDeliveryConfig != nil {
	// 	err = db_organization.Query.UpdateOrganizationEggSalesConfig(ctx, db_organization.UpdateOrganizationEggSalesConfigParams{
	// 		OrganizationID: organization.ID,
	// 		// PreferredUnit: nil,
	// 		EggPriceMin:   req.Msg.EggDeliveryConfig.MinimumJobCost,
	// 		EggPrice:      req.Msg.EggDeliveryConfig.EggCost,
	// 		EggPricePerKm: req.Msg.EggDeliveryConfig.DutyCostPerKm,
	// 		IsEnabled: pgtype.Bool{
	// 			Bool:  req.Msg.EggDeliveryConfig.IsEnabled,
	// 			Valid: true,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, connect.NewError(
	// 			connect.CodeInternal,
	// 			errors.New("internal server error"),
	// 		)
	// 	}
	// 	configSet = true
	// }

	// if req.Msg.CoconutDeliveryConfig != nil {
	// 	err = db_organization.Query.UpdateOrganizationCoconutSalesConfig(ctx, db_organization.UpdateOrganizationCoconutSalesConfigParams{
	// 		OrganizationID: organization.ID,
	// 		// PreferredUnit: nil,
	// 		CoconutPriceMin:   req.Msg.CoconutDeliveryConfig.MinimumJobCost,
	// 		CoconutPrice:      req.Msg.CoconutDeliveryConfig.CoconutCost,
	// 		CoconutPricePerKm: req.Msg.CoconutDeliveryConfig.DutyCostPerKm,
	// 		IsEnabled: pgtype.Bool{
	// 			Bool:  req.Msg.CoconutDeliveryConfig.IsEnabled,
	// 			Valid: true,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, connect.NewError(
	// 			connect.CodeInternal,
	// 			errors.New("internal server error"),
	// 		)
	// 	}
	// 	configSet = true
	// }

	if req.Msg.SubmitForApproval != nil && *req.Msg.SubmitForApproval {
		if configSet == false {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("config not set"),
			)
		}
		// TODO: Verify all required fields have been provided
	}

	return connect.NewResponse(&endpoints.CreateOrganizationResponse{
		OrganizationId: organization.ID,
	}), nil
}
