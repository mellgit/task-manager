package task

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mellgit/task-manager/internal/middleware"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	GroupHandler(app *fiber.App)
}

type handler struct {
	service Service
	Logger  *log.Entry
}

func NewHandler(service Service, logger *log.Entry) Handler {
	return &handler{
		service: service,
		Logger:  logger,
	}
}

func (h *handler) GroupHandler(app *fiber.App) {
	group := app.Group("/api/tasks", middleware.JWTProtected())
	group.Post("/", h.Create)

}

func (h *handler) Create(ctx *fiber.Ctx) error {

	payload := &TaskRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	userID := ctx.Locals("user_id").(string)

	if err := h.service.CreateTask(uuid.MustParse(userID), payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "CreateTask",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}

	return ctx.SendStatus(fiber.StatusCreated)

}
