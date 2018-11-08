package handlers

import (
	"context"
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/client-go/rest"
)

func NewHandler(cfg *rest.Config) sdk.Handler {
	return &Handler{
		cfg: cfg,
	}
}

type Handler struct {
	cfg *rest.Config
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *integreatly.ManagedServiceNamespace:
		return handleManagedServiceNamespace(ctx, event, o, h.cfg)
	}

	return nil
}
