package v1

import (
	"context"
	"errors"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/service"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"github.com/keenywheels/backend/pkg/httputils"
)

// SearchTokenInfo searches for token information based on the provided parameters.
func (c *Controller) SearchTokenInfo(
	ctx context.Context,
	params gen.SearchTokenInfoParams,
) (gen.SearchTokenInfoRes, error) {
	var (
		op  = "Controller.SearchTokenInfo"
		log = ctxutils.GetLogger(ctx)
	)

	tokensInfo, err := c.svc.SearchTokenInfo(ctx, params.Token)
	if err != nil {
		log.Errorf("[%s] failed to search token info: %v", op, err)

		switch {
		case errors.Is(err, service.ErrNotFound):
			return &gen.SearchTokenInfoNotFound{
				Error: httputils.ErrorNotFound,
			}, nil
		}

		return &gen.SearchTokenInfoInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	resp := gen.SearchTokenInfoOKApplicationJSON(convertToSearchTokenInfoResp(tokensInfo))

	return &resp, nil
}
