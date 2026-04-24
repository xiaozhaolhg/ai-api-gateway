package handler

import "context"

type Model struct {
	ID       string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ModelsListResp struct {
	Object string `json:"object"`
	Data   []Model `json:"data"`
}

type ModelsService interface {
	ListModels(ctx context.Context) (*ModelsListResp, error)
}

type ModelsHandler struct {
	svc ModelsService
}

func NewModelsHandler(svc ModelsService) *ModelsHandler {
	return &ModelsHandler{svc: svc}
}

func (h *ModelsHandler) ListModels() (*ModelsListResp, error) {
	return h.svc.ListModels(context.Background())
}