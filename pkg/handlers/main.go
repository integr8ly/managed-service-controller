package handlers

import (
	"context"
	integreatly "github.com/integr8ly/managed-service-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/client-go/rest"
)

type Handler struct {
	msnHandler *ManagedServiceNamespaceHandler
}

func NewHandler(cfg *rest.Config) sdk.Handler {
	return &Handler{
		msnHandler: NewManagedServiceNamespaceHandler(cfg),
	}
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *integreatly.ManagedServiceNamespace:
		return h.msnHandler.Handle(ctx, event, o)
	}

	return nil
}
