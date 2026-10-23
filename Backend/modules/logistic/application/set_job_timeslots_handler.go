package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) SetJobTimeslotsHandler(
	ctx context.Context,
	req *connect.Request[endpoints.SetJobTimeslotCommand],
) (*connect.Response[endpoints.SetJobTimeslotResponse], error) {
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

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: vehicle.OrganizationID,
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
			errors.New("not permitted"),
		)
	}

	var weekdays []int32

	existingTimeslots, err := db_logistic.Query.GetTimeslots(ctx, req.Msg.VehicleId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	for _, jobTimeslot := range req.Msg.JobTimeslots {
		for _, weekday := range jobTimeslot.Weekdays {
			weekdays = append(weekdays, int32(weekday.Number()))
		}

		StartHH := *jobTimeslot.StartHH
		StartMM := *jobTimeslot.StartMM
		EndHH := *jobTimeslot.EndHH
		EndMM := *jobTimeslot.EndMM

		_, err := db_logistic.Query.CreateTimeslot(ctx, db_logistic.CreateTimeslotParams{
			VehicleID: vehicle.ID,
			StartHh:   int16(StartHH),
			StartMm:   int16(StartMM),
			EndHh:     int16(EndHH),
			EndMm:     int16(EndMM),
			Wkdays:    weekdays,
		})

		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
	}

	for _, jobTimeslot := range existingTimeslots {
		db_logistic.Query.DeleteTimeslot(ctx, jobTimeslot.ID)
	}

	return connect.NewResponse(&endpoints.SetJobTimeslotResponse{}), nil
}
