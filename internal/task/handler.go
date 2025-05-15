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
	group.Get("/protected", h.Protected)
	group.Get("/:task_id", h.GetTask)
	group.Delete("/:task_id", h.DeleteTask)
	group.Patch("/:task_id", h.UpdateTask)

}

// Create
// @Summary      Create
// @Description  Create task
// @Security ApiKeyAuth
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param 		 request body TaskRequest true "body"
// @Success      200 {object} int
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/tasks [post]
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

// List
// @Summary      List
// @Description  List tasks
// @Security ApiKeyAuth
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Success      200 {array} Task
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api [get]
func (h *Handler) List(ctx *fiber.Ctx) error {

	userID := ctx.Locals("user_id").(string)
	tasks, err := h.service.ListTasks(uuid.MustParse(userID))
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "List",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}

	return ctx.Status(fiber.StatusOK).JSON(tasks)
}

// GetTask
// @Summary      GetTask
// @Description  Get one task
// @Security ApiKeyAuth
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        task_id path string true "task_id"
// @Success      200 {object} Task
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/tasks/{task_id} [get]
func (h *Handler) GetTask(ctx *fiber.Ctx) error {

	userID := ctx.Locals("user_id").(string)
	taskID := ctx.Params("task_id")

	task, err := h.service.GetTask(uuid.MustParse(userID), uuid.MustParse(taskID))
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "GetTask",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}
	return ctx.Status(fiber.StatusOK).JSON(task)

}

// DeleteTask
// @Summary      DeleteTask
// @Description  Delete task
// @Security ApiKeyAuth
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        task_id path string true "task_id"
// @Success      200 {object} int
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/tasks/{task_id} [delete]
func (h *Handler) DeleteTask(ctx *fiber.Ctx) error {

	userID := ctx.Locals("user_id").(string)
	taskID := ctx.Params("task_id")
	if err := h.service.DeleteTask(uuid.MustParse(userID), uuid.MustParse(taskID)); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "DeleteTask",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

// UpdateTask
// @Summary      UpdateTask
// @Description  Update task
// @Security ApiKeyAuth
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param 		 request body TaskRequest true "body"
// @Success      200 {object} int
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/tasks/{task_id} [patch]
func (h *Handler) UpdateTask(ctx *fiber.Ctx) error {

	userID := ctx.Locals("user_id").(string)
	payload := &TaskRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}
	if err := h.service.UpdateTask(uuid.MustParse(userID), payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "DeleteTask",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// Protected
// @Summary      Protected
// @Description  Protected check authorized user
// @Security ApiKeyAuth
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/tasks/protected [get]
func (h *Handler) Protected(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(MessageResponse{Message: "authorized user"})
}
