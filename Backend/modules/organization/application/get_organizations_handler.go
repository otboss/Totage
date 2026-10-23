package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) GetOrganizationsHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetOrganizationsQuery],
) (*connect.Response[endpoints.GetOrganizationsResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	var offsetId int64 = 0
	if req.Msg.Pagination.OffsetId != nil {
		offsetId = *req.Msg.Pagination.OffsetId
	}

	filterText := ""
	if req.Msg.Pagination.FilterText != nil {
		filterText = *req.Msg.Pagination.FilterText
	}

	organizations, err := db_organization.Query.GetOrganizationsForUser(ctx, db_organization.GetOrganizationsForUserParams{
		OffsetID:   offsetId,
		Limitation: req.Msg.Pagination.Limit,
		FilterText: filterText,
		UserID:     claims.UserID,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var responseBody []*endpoints.GetOrganizationsResponse_Organization

	for _, organization := range organizations {
		var avatarUrl *string
		if organization.AvatarUrl.Valid == true {
			avatarUrl = &organization.AvatarUrl.String
		}
		memberCount, err := db_organization.Query.GetAccountsCount(ctx, organization.ID)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
		responseBody = append(responseBody, &endpoints.GetOrganizationsResponse_Organization{
			OrganizationId:   organization.ID,
			OrganizationName: organization.OrganizationName,
			AvatarUrl:        avatarUrl,
			Rating:           int64(organization.Rating),
			MemberCount:      memberCount,
			IsApproved:       organization.IsApproved,
		})
	}

	return connect.NewResponse(&endpoints.GetOrganizationsResponse{
		Organizations: responseBody,
	}), nil
}
