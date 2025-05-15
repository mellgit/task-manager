package task

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mellgit/task-manager/internal/middleware"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	service Service
	Logger  *log.Entry
}

func NewHandler(service Service, logger *log.Entry) *Handler {
	return &Handler{
		service: service,
		Logger:  logger,
	}
}

func (h *Handler) GroupHandler(app *fiber.App) {
	group := app.Group("/api/tasks", middleware.JWTProtected())
	group.Post("/", h.Create)
	group.Get("/", h.List)

}

func (h *Handler) Create(ctx *fiber.Ctx) error {

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

func (h *Handler) List(ctx *fiber.Ctx) error {

	userID := ctx.Locals("user_id").(string)
	tasks, err := h.service.ListTasks(uuid.MustParse(userID))
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "CreateTask",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}

	return ctx.Status(fiber.StatusOK).JSON(tasks)
}
