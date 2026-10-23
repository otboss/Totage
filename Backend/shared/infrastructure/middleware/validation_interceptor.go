package middleware

import (
	"context"

	"connectrpc.com/connect"
)

func ValidationInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {

			// Check if request implements Validate()
			if validator, ok := req.Any().(interface{ Validate() error }); ok {
				if err := validator.Validate(); err != nil {
					return nil, connect.NewError(
						connect.CodeInvalidArgument,
						err,
					)
				}
			}

			return next(ctx, req)
		}
	}
}
