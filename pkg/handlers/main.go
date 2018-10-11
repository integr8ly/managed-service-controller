package handlers

import (
	"context"
	integreatly "github.com/integr8ly/managed-services-controller/pkg/apis/integreatly/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/client-go/kubernetes"
)

func NewHandler(client kubernetes.Interface) sdk.Handler {
	return &Handler{
	    client: client,
	}
}

type Handler struct {
	client kubernetes.Interface
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *integreatly.ManagedServiceNamespace:
		return handleManagedServiceNamespace(ctx, event, o, h.client)
	}

	return nil
}