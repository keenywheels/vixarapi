package v1

import (
	"context"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"github.com/keenywheels/backend/pkg/httputils"
)

// GetAllInterest get interest for specified token in all time
func (c *Controller) GetAllInterest(
	ctx context.Context,
	req *gen.GetAllInterestRequest,
) (gen.GetAllInterestRes, error) {
	op := "Controller.GetAllInterest"
	log := ctxutils.GetLogger(ctx)

	if err := req.Validate(); err != nil {
		log.Errorf("[%s] failed to validate request: %v", op, err)

		return &gen.GetAllInterestBadRequest{
			Error: httputils.ErrorBadRequest,
		}, nil
	}

	interests, err := c.svc.GetAllInterest(ctx, req.Token)
	if err != nil {
		log.Errorf("[%s] failed to get interest: %v", op, err)

		return &gen.GetAllInterestInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	resp := gen.GetAllInterestOKApplicationJSON(convertToGetAllInterestResp(interests))

	return &resp, nil
}
