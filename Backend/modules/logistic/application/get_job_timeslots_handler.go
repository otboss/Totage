package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) GetJobTimeslotsHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetJobTimeslotsQuery],
) (*connect.Response[endpoints.GetJobTimeslotsResponse], error) {

	jobTimeslots, err := db_logistic.Query.GetTimeslots(ctx, req.Msg.VehicleId)
	if err == pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodeNotFound,
			errors.New("not found"),
		)
	}
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var response endpoints.GetJobTimeslotsResponse

	for _, jobTimeslot := range jobTimeslots {
		var timeslot common.JobTimeslot

		for _, weekdays := range jobTimeslot.Weekdays {
			timeslot.Weekdays = append(timeslot.Weekdays, common.Weekdays(weekdays))
		}

		StartHh := int32(jobTimeslot.StartHh)
		StartMm := int32(jobTimeslot.StartMm)
		EndHh := int32(jobTimeslot.EndHh)
		EndMm := int32(jobTimeslot.EndMm)

		timeslot.StartHH = &StartHh
		timeslot.StartMM = &StartMm
		timeslot.EndHH = &EndHh
		timeslot.EndMM = &EndMm

		response.JobTimeslot = append(response.JobTimeslot, &timeslot)
	}

	return connect.NewResponse(&response), nil
}
