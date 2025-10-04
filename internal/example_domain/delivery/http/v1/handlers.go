package v1

import (
	"context"
	"fmt"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"github.com/keenywheels/backend/pkg/httputils"
)

// Example ...
func (c *Controller) Example(ctx context.Context, req *gen.ExampleDomainRequest) (gen.ExampleRes, error) {
	op := "Controller.Example"
	log := ctxutils.GetLogger(ctx)

	if err := req.Validate(); err != nil {
		log.Errorf("[%s] failed to validate request: %v", op, err)

		return &gen.ExampleBadRequest{
			Error: httputils.ErrorBadRequest,
		}, nil
	}

	if req.Username != "ilya228" {
		log.Errorf("[%s] got wrong creds for %s", op, req.Username)

		return &gen.ExampleUnauthorized{
			Error: httputils.ErrorUnathorized,
		}, nil
	}

	log.Infof("[%s] got request: %+v", op, req)

	return &gen.ExampleDomainResponse{
		Message: fmt.Sprintf("Hello, %s", op, req.Username),
		Length:  len(req.Username),
	}, nil
}
