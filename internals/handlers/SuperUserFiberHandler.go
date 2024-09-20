package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/responses"
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
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid input", nil, err.Error()))
	}

	createdSuperUser, err := h.service.CreateSuperUser(context.Background(), &superUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to create SuperUser", nil, err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(responses.NewFiberResponse(c, fiber.StatusCreated, "SuperUser created successfully", createdSuperUser, nil))
}

// GetAllSuperUsersHandler retrieves all SuperUsers and returns them in a Fiber response
func (h *SuperUserFiberHandler) GetAllSuperUsersHandler(c *fiber.Ctx) error {
	superUsers, err := h.service.GetAllSuperUsers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to retrieve SuperUsers", nil, err.Error()))
	}

	if len(superUsers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewFiberResponse(c, fiber.StatusNotFound, "No SuperUsers found", nil, nil))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "SuperUsers retrieved successfully", superUsers, nil))
}

// Get SuperUser by ID
func (h *SuperUserFiberHandler) GetSuperUserByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	superUser, err := h.service.GetSuperUserByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewFiberResponse(c, fiber.StatusNotFound, "SuperUser not found", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "SuperUser retrieved successfully", superUser, nil))
}

// Get SuperUser by email
func (h *SuperUserFiberHandler) GetSuperUserByEmailHandler(c *fiber.Ctx) error {
	email := c.Params("email")

	superUser, err := h.service.GetSuperUserByEmail(context.Background(), email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewFiberResponse(c, fiber.StatusNotFound, "SuperUser not found", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "SuperUser retrieved successfully", superUser, nil))
}

// Get SuperUser by username
func (h *SuperUserFiberHandler) GetSuperUserByUsernameHandler(c *fiber.Ctx) error {
	username := c.Params("username")

	superUser, err := h.service.GetSuperUserByUsername(context.Background(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewFiberResponse(c, fiber.StatusNotFound, "SuperUser not found", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "SuperUser retrieved successfully", superUser, nil))
}

// Enable 2FA for SuperUser
func (h *SuperUserFiberHandler) Enable2FAForSuperUserHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	var body struct {
		Secret string `json:"secret"`
	}
	if err := c.BodyParser(&body); err != nil || body.Secret == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid input", nil, "Secret is missing or invalid"))
	}

	if err := h.service.Enable2FAForSuperUser(context.Background(), id, body.Secret); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to enable 2FA", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "2FA enabled", nil, nil))
}

// Disable 2FA for SuperUser
func (h *SuperUserFiberHandler) Disable2FAForSuperUserHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	if err := h.service.Disable2FAForSuperUser(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to disable 2FA", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "2FA disabled", nil, nil))
}

// Get all 2FA-enabled SuperUsers
func (h *SuperUserFiberHandler) GetAll2FAEnabledSuperUsersHandler(c *fiber.Ctx) error {
	superUsers, err := h.service.GetAll2FAEnabledSuperUsers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to retrieve 2FA-enabled SuperUsers", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "2FA-enabled SuperUsers retrieved", superUsers, nil))
}

// Update SuperUser role
func (h *SuperUserFiberHandler) UpdateSuperUserRoleHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil || body.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid input", nil, "Role is missing or invalid"))
	}

	if err := h.service.UpdateSuperUserRole(context.Background(), id, body.Role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to update role", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "Role updated", nil, nil))
}

// Update SuperUser permissions
func (h *SuperUserFiberHandler) UpdateSuperUserPermissionsHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	var body struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.Permissions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid input", nil, "Permissions are missing or invalid"))
	}

	if err := h.service.UpdateSuperUserPermissions(context.Background(), id, body.Permissions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to update permissions", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "Permissions updated", nil, nil))
}

// Update specific SuperUser field
func (h *SuperUserFiberHandler) UpdateSuperUserFieldHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	field := c.Params("field")
	var value interface{}
	if err := c.BodyParser(&value); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid input", nil, "Field value is invalid"))
	}

	if err := h.service.UpdateSuperUserField(context.Background(), id, field, value); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to update field", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "Field updated", nil, nil))
}

// Generate and set reset token
func (h *SuperUserFiberHandler) GenerateAndSetResetTokenHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	token, err := h.service.GenerateAndSetResetToken(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to generate reset token", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "Reset token generated", map[string]string{"reset_token": token}, nil))
}

// Clear reset token
func (h *SuperUserFiberHandler) ClearResetTokenHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	if err := h.service.ClearResetToken(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to clear reset token", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "Reset token cleared", nil, nil))
}

// Delete SuperUser by ID
func (h *SuperUserFiberHandler) DeleteSuperUserByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil, err.Error()))
	}

	if err := h.service.DeleteSuperUserByID(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to delete SuperUser", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "SuperUser deleted", nil, nil))
}

// Search SuperUsers
func (h *SuperUserFiberHandler) SearchSuperUsersHandler(c *fiber.Ctx) error {
	searchQuery := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "created_at")

	superUsers, err := h.service.SearchSuperUsers(context.Background(), searchQuery, page, limit, sortBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to search SuperUsers", nil, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewFiberResponse(c, fiber.StatusOK, "SuperUsers search successful", superUsers, nil))
}
