package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/services"
	"github.com/lordofthemind/EventifyGo/internals/types"
)

type SuperUserFiberHandler struct {
	service services.SuperUserServiceInterface
}

func NewSuperUserFiberHandler(service services.SuperUserServiceInterface) *SuperUserFiberHandler {
	return &SuperUserFiberHandler{service: service}
}

// Create SuperUser handler
func (h *SuperUserFiberHandler) CreateSuperUserHandler(c *fiber.Ctx) error {
	var superUser types.SuperUserType
	if err := c.BodyParser(&superUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	createdSuperUser, err := h.service.CreateSuperUser(context.Background(), &superUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(createdSuperUser)
}

// Get SuperUser by ID
func (h *SuperUserFiberHandler) GetSuperUserByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	superUser, err := h.service.GetSuperUserByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(superUser)
}

// Get SuperUser by email
func (h *SuperUserFiberHandler) GetSuperUserByEmailHandler(c *fiber.Ctx) error {
	email := c.Params("email")

	superUser, err := h.service.GetSuperUserByEmail(context.Background(), email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(superUser)
}

// Get SuperUser by username
func (h *SuperUserFiberHandler) GetSuperUserByUsernameHandler(c *fiber.Ctx) error {
	username := c.Params("username")

	superUser, err := h.service.GetSuperUserByUsername(context.Background(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(superUser)
}

// Enable 2FA for SuperUser
func (h *SuperUserFiberHandler) Enable2FAForSuperUserHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var body struct {
		Secret string `json:"secret"`
	}
	if err := c.BodyParser(&body); err != nil || body.Secret == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.Enable2FAForSuperUser(context.Background(), id, body.Secret); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "2FA enabled"})
}

// Disable 2FA for SuperUser
func (h *SuperUserFiberHandler) Disable2FAForSuperUserHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := h.service.Disable2FAForSuperUser(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "2FA disabled"})
}

// Get all 2FA-enabled SuperUsers
func (h *SuperUserFiberHandler) GetAll2FAEnabledSuperUsersHandler(c *fiber.Ctx) error {
	superUsers, err := h.service.GetAll2FAEnabledSuperUsers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(superUsers)
}

// Update SuperUser role
func (h *SuperUserFiberHandler) UpdateSuperUserRoleHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil || body.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.UpdateSuperUserRole(context.Background(), id, body.Role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Role updated"})
}

// Update SuperUser permissions
func (h *SuperUserFiberHandler) UpdateSuperUserPermissionsHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var body struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.Permissions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.UpdateSuperUserPermissions(context.Background(), id, body.Permissions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Permissions updated"})
}

// Update specific SuperUser field
func (h *SuperUserFiberHandler) UpdateSuperUserFieldHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	field := c.Params("field")
	var value interface{}
	if err := c.BodyParser(&value); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.UpdateSuperUserField(context.Background(), id, field, value); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Field updated"})
}

// Generate and set reset token
func (h *SuperUserFiberHandler) GenerateAndSetResetTokenHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	token, err := h.service.GenerateAndSetResetToken(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"reset_token": token})
}

// Clear reset token
func (h *SuperUserFiberHandler) ClearResetTokenHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := h.service.ClearResetToken(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Reset token cleared"})
}

// Delete SuperUser by ID
func (h *SuperUserFiberHandler) DeleteSuperUserByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := h.service.DeleteSuperUserByID(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "SuperUser deleted"})
}

// Search SuperUsers
func (h *SuperUserFiberHandler) SearchSuperUsersHandler(c *fiber.Ctx) error {
	searchQuery := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "created_at")

	superUsers, err := h.service.SearchSuperUsers(context.Background(), searchQuery, page, limit, sortBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(superUsers)
}
