package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) ApproveOrganizationHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ApproveOrganizationCommand],
) (*connect.Response[endpoints.ApproveOrganizationResponse], error) {
	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	sysAdmin, err := db_user.Query.GetSystemAdministratorByUserId(ctx, claims.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("not authorized to perform action")
		}
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, err = db_organization.Query.CreateOrganizationApproval(ctx, db_organization.CreateOrganizationApprovalParams{
		OrganizationID:        req.Msg.OrganizationId,
		SystemAdministratorID: sysAdmin.ID,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	approvals, err := db_organization.Query.GetOrganizationApprovals(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	systemAdmins, err := db_user.Query.GetSystemAdministrators(ctx)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if len(approvals)/len(systemAdmins)*100 >= 51 {
		err = db_organization.Query.UpdateOrganizationApproval(ctx, req.Msg.OrganizationId)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
	}

	return connect.NewResponse(&endpoints.ApproveOrganizationResponse{}), nil
}
