package gateway

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"super-node-core/internal/cocoon"
)

type Handler struct {
	Orchestrator *cocoon.Orchestrator
}

func NewHandler(orch *cocoon.Orchestrator) *Handler {
	return &Handler{Orchestrator: orch}
}

type TransactionRequest struct {
	UserID       string  `json:"user_id"`
	Amount       float64 `json:"amount"`
	SourceWallet string  `json:"source_wallet"`
}

func (h *Handler) SubmitTransaction(c *fiber.Ctx) error {
	var req TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Create internal transaction object
	tx := &cocoon.Transaction{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		Amount:       req.Amount,
		SourceWallet: req.SourceWallet,
		Status:       "PENDING",
	}

	// In a real event-driven architecture, we might just push to Redpanda here and return "Accepted".
	// For this synchronous demo requirement, we wait for the result.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.Orchestrator.ExecuteWorkflow(ctx, tx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":        "Transaction processed",
		"transaction_id": tx.ID,
		"status":         tx.Status,
	})
}
