package application

import (
	"context"
	"errors"
	"time"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"

	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	// db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handlers) GetLogisticScheduleHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetLogisticScheduleQuery],
) (*connect.Response[endpoints.GetLogisticScheduleResponse], error) {

	organization, err := db_organization.Query.GetOrganization(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	jobMonth := pgtype.Int2{
		Int16: int16(req.Msg.Month),
		Valid: true,
	}

	jobYear := pgtype.Int2{
		Int16: int16(time.Now().Year()),
		Valid: true,
	}

	scheduledOrders, err := db_order.Query.GetOrdersIncompletedForMonth(ctx, db_order.GetOrdersIncompletedForMonthParams{
		OrganizationID: organization.ID,
		JobMonth:       jobMonth,
		JobYear:        jobYear,
	})

	var schedules []*endpoints.GetLogisticScheduleResponse_Schedule

	for _, order := range scheduledOrders {
		var scheduleTruck *common.Truck

		if order.Type == db_order.SharedOrderTypeFREIGHT {
			fullOrder, err := db_order.Query.GetOrder(ctx, order.ID)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInternal,
					errors.New("internal server error"),
				)
			}

			freightOrder, err := db_order.Query.GetFreightOrderByOrderId(ctx, pgtype.UUID{
				Bytes: fullOrder.ID.Bytes,
				Valid: true,
			})
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInternal,
					errors.New("internal server error"),
				)
			}

			truck, err := db_logistic.Query.GetTruck(ctx, freightOrder.TruckID)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInternal,
					errors.New("internal server error"),
				)
			}

			var truckBedHeight *float32
			truckBedHeightVal := float32(truck.BedHeightMm.Float64)
			if truck.BedHeightMm.Valid {
				truckBedHeight = &truckBedHeightVal
			}

			bedSize := common.BedSize{
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

			// truckProperty, err := local_utils.TruckPropertyMap.DbTruckPropertyToProtoTruckProperty(truck.Property)
			// if err != nil {
			// 	return nil, connect.NewError(
			// 		connect.CodeInternal,
			// 		errors.New("internal server error"),
			// 	)
			// }

			scheduleTruck = &common.Truck{
				BedSize:           &bedSize,
				Property:          common.FreightTruckProperty(truckPropertyNum),
				FreightPricePerKm: float32(truck.FreightPricePerKm),
				Cost:              float32(truck.Cost),
			}
		}

		scheduleItem := endpoints.GetLogisticScheduleResponse_Schedule_ScheduleItem{
			OrderId:                  order.ID.String(),
			JobRoundTripDistanceInKm: order.DistanceRoundTripInKm,
			OrderType:                string(order.Type),
			Truck:                    scheduleTruck,
		}

		schedule := endpoints.GetLogisticScheduleResponse_Schedule{
			ScheduleItem: &scheduleItem,
			Timestamp:    timestamppb.Now(),
		}

		schedules = append(schedules, &schedule)
	}

	res := connect.NewResponse(&endpoints.GetLogisticScheduleResponse{
		Schedule: schedules,
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}
