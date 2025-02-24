package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-go-common-shared/kafka"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	"github.com/ngdangkietswe/swe-integration-service/data/repository/cdcauthusers"
	"go.uber.org/zap"
)

type CdcAuthUsersHandler struct {
	logger           *logger.Logger
	cdcAuthUsersRepo cdcauthusers.ICdcAuthUsersRepository
}

func NewCdcAuthUsersHandler(logger *logger.Logger, cdcAuthUsersRepo cdcauthusers.ICdcAuthUsersRepository) *CdcAuthUsersHandler {
	return &CdcAuthUsersHandler{
		logger:           logger,
		cdcAuthUsersRepo: cdcAuthUsersRepo,
	}
}

// Handle handles the CDC event for the CdcAuthUsers table.
func (h *CdcAuthUsersHandler) Handle(cdcEventMsg *kafka.CdcEventMsg) {
	if cdcEventMsg.Op == kafka.CdcOperationDelete {
		h.handleDelete(cdcEventMsg)
	} else {
		h.handleUpsert(cdcEventMsg)
	}
}

// handleUpsert handles the CDC event for inserting a new record.
func (h *CdcAuthUsersHandler) handleUpsert(cdcEventMsg *kafka.CdcEventMsg) {
	if err := h.cdcAuthUsersRepo.UpsertByCdcEventMsg(context.Background(), cdcEventMsg); err != nil {
		h.logger.Error("Failed to upsert record by CDC event message", zap.String("table", "CdcAuthUsers"), zap.Error(err))
	}
}

// handleDelete handles the CDC event for deleting an existing record.
func (h *CdcAuthUsersHandler) handleDelete(cdcEventMsg *kafka.CdcEventMsg) {
	if err := h.cdcAuthUsersRepo.DeleteById(context.Background(), uuid.MustParse(cdcEventMsg.Before["id"].(string))); err != nil {
		h.logger.Error("Failed to delete record by CDC event message", zap.String("table", "CdcAuthUsers"), zap.Error(err))
	}
}
