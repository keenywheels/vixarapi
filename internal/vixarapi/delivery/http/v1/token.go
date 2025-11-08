package v1

import (
	"context"
	"errors"
	"time"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/service"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"github.com/keenywheels/backend/pkg/httputils"
)

// SearchTokenInfo searches for token information based on the provided parameters.
func (c *Controller) SearchTokenInfo(
	ctx context.Context,
	req *gen.SearchTokenInfoRequest,
) (gen.SearchTokenInfoRes, error) {
	var (
		op  = "Controller.SearchTokenInfo"
		log = ctxutils.GetLogger(ctx)
	)

	end := time.Now().UTC()
	if req.End.Set {
		end = req.End.Value.UTC()
	}

	tokensInfo, err := c.svc.SearchTokenInfo(ctx, &service.SearchTokenInfoParams{
		Token: req.Token,
		Start: req.Start.UTC(),
		End:   end,
	})
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
