package application

// import (
// 	"context"
// 	"errors"
// 	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
// 	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
// 	common "waterordersystem/GOBACKEND/shared/api"

// 	"connectrpc.com/connect"
// 	"github.com/jackc/pgx/v5"
// )

// func (s *Handlers) GetSalesUnitsHandler(
// 	ctx context.Context,
// 	req *connect.Request[endpoints.GetSalesUnitsQuery],
// ) (*connect.Response[endpoints.GetSalesUnitsResponse], error) {

// 	var unit string
// 	var unitCost float64
// 	var quantity []int64
// 	var salesUnits []db_organization.GetSalesUnitsRow

// 	switch req.Msg.OrderType {
// 	case common.OrderType_WATER_PURIFIED:
// 		waterSalesConfig, err := db_organization.Query.GetPotableWaterSalesConfig(ctx, req.Msg.OrganizationId)
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				return nil, errors.New("order type not supported for organization")
// 			}
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 		unitCost = waterSalesConfig.WaterPricePerGallon.Float64
// 		unit = string(waterSalesConfig.PreferredUnit)

// 		// TODO: Do unit conversion if unit is not the default
// 		salesUnits, err = db_organization.Query.GetSalesUnits(ctx, db_organization.GetSalesUnitsParams{
// 			OrganizationID: req.Msg.OrganizationId,
// 			OrderType:      db_organization.SharedOrderTypeWATERPURIFIED,
// 		})
// 		if err != nil {
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 	case common.OrderType_WATER_UNPURIFIED:
// 		waterSalesConfig, err := db_organization.Query.GetIrrigationWaterSalesConfig(ctx, req.Msg.OrganizationId)
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				return nil, errors.New("order type not supported for organization")
// 			}
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 		unitCost = waterSalesConfig.WaterPricePerGallon.Float64
// 		unit = string(waterSalesConfig.PreferredUnit)

// 		// TODO: Do unit conversion if unit is not the default

// 		salesUnits, err = db_organization.Query.GetSalesUnits(ctx, db_organization.GetSalesUnitsParams{
// 			OrganizationID: req.Msg.OrganizationId,
// 			OrderType:      db_organization.SharedOrderTypeWATERUNPURIFIED,
// 		})
// 		if err != nil {
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 	case common.OrderType_BLOCK:
// 		blockSalesConfig, err := db_organization.Query.GetBlockSalesConfig(ctx, req.Msg.OrganizationId)
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				return nil, errors.New("order type not supported for organization")
// 			}
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 		unitCost = blockSalesConfig.BlockPrice
// 		unit = string(blockSalesConfig.PreferredUnit)

// 		salesUnits, err = db_organization.Query.GetSalesUnits(ctx, db_organization.GetSalesUnitsParams{
// 			OrganizationID: req.Msg.OrganizationId,
// 			OrderType:      db_organization.SharedOrderTypeBLOCK,
// 		})
// 		if err != nil {
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 	case common.OrderType_COCONUT:
// 		coconutSalesConfig, err := db_organization.Query.GetCoconutSalesConfig(ctx, req.Msg.OrganizationId)
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				return nil, errors.New("order type not supported for organization")
// 			}
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 		unitCost = coconutSalesConfig.CoconutPrice
// 		unit = string(coconutSalesConfig.PreferredUnit)

// 		salesUnits, err = db_organization.Query.GetSalesUnits(ctx, db_organization.GetSalesUnitsParams{
// 			OrganizationID: req.Msg.OrganizationId,
// 			OrderType:      db_organization.SharedOrderTypeCOCONUT,
// 		})
// 		if err != nil {
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 	case common.OrderType_EGG:
// 		eggSalesConfig, err := db_organization.Query.GetEggSalesConfig(ctx, req.Msg.OrganizationId)
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				return nil, errors.New("order type not supported for organization")
// 			}
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 		unitCost = eggSalesConfig.EggPrice
// 		unit = string(eggSalesConfig.PreferredUnit)

// 		salesUnits, err = db_organization.Query.GetSalesUnits(ctx, db_organization.GetSalesUnitsParams{
// 			OrganizationID: req.Msg.OrganizationId,
// 			OrderType:      db_organization.SharedOrderTypeEGG,
// 		})
// 		if err != nil {
// 			return nil, connect.NewError(
// 				connect.CodeInternal,
// 				errors.New("internal server error"),
// 			)
// 		}

// 	default:
// 		return nil, errors.New("order type not supported")
// 	}

// 	for _, saleUnit := range salesUnits {
// 		quantity = append(quantity, int64(saleUnit.Quantity))
// 	}

// 	return connect.NewResponse(&endpoints.GetSalesUnitsResponse{
// 		UnitCost:   unitCost,
// 		Unit:       unit,
// 		Quantities: quantity,
// 	}), nil
// }
