package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/responses" // Import the responses package
	"github.com/lordofthemind/EventifyGo/internals/services"
	"github.com/lordofthemind/EventifyGo/internals/types"
)

type SuperUserGinHandler struct {
	service services.SuperUserServiceInterface
}

func NewSuperUserGinHandler(service services.SuperUserServiceInterface) *SuperUserGinHandler {
	return &SuperUserGinHandler{service: service}
}

// Create SuperUser handler
func (h *SuperUserGinHandler) CreateSuperUserHandler(c *gin.Context) {
	var superUser types.SuperUserType
	if err := c.ShouldBindJSON(&superUser); err != nil {
		// Use standardized response for invalid input
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	createdSuperUser, err := h.service.CreateSuperUser(c.Request.Context(), &superUser)
	if err != nil {
		// Use standardized response for internal server errors
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to create SuperUser", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Use standardized response for successful creation
	response := responses.NewGinResponse(c, http.StatusCreated, "SuperUser created successfully", createdSuperUser, nil)
	c.JSON(http.StatusCreated, response)
}

// GetAllSuperUsersHandler retrieves all SuperUsers and returns them in a standardized response
func (h *SuperUserGinHandler) GetAllSuperUsersHandler(c *gin.Context) {
	superUsers, err := h.service.GetAllSuperUsers(c.Request.Context())
	if err != nil {
		// Use standardized response for internal server errors
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to retrieve SuperUsers", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// If no SuperUsers are found, return a standardized response for empty result
	if len(superUsers) == 0 {
		response := responses.NewGinResponse(c, http.StatusNotFound, "No SuperUsers found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	// Use standardized response for successful retrieval
	response := responses.NewGinResponse(c, http.StatusOK, "SuperUsers retrieved successfully", superUsers, nil)
	c.JSON(http.StatusOK, response)
}

// Get SuperUser by ID handler
func (h *SuperUserGinHandler) GetSuperUserByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		// Use standardized response for invalid ID format
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	superUser, err := h.service.GetSuperUserByID(c.Request.Context(), id)
	if err != nil {
		// Use standardized response for not found error
		response := responses.NewGinResponse(c, http.StatusNotFound, "SuperUser not found", nil, err.Error())
		c.JSON(http.StatusNotFound, response)
		return
	}

	// Use standardized response for successful retrieval
	response := responses.NewGinResponse(c, http.StatusOK, "SuperUser retrieved successfully", superUser, nil)
	c.JSON(http.StatusOK, response)
}

// Get SuperUser by email
func (h *SuperUserGinHandler) GetSuperUserByEmailHandler(c *gin.Context) {
	email := c.Param("email")

	superUser, err := h.service.GetSuperUserByEmail(c.Request.Context(), email)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusNotFound, "SuperUser not found", nil, err.Error())
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "SuperUser retrieved successfully", superUser, nil)
	c.JSON(http.StatusOK, response)
}

// Get SuperUser by username
func (h *SuperUserGinHandler) GetSuperUserByUsernameHandler(c *gin.Context) {
	username := c.Param("username")

	superUser, err := h.service.GetSuperUserByUsername(c.Request.Context(), username)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusNotFound, "SuperUser not found", nil, err.Error())
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "SuperUser retrieved successfully", superUser, nil)
	c.JSON(http.StatusOK, response)
}

// Enable 2FA for SuperUser
func (h *SuperUserGinHandler) Enable2FAForSuperUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var body struct {
		Secret string `json:"secret"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Secret == "" {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.service.Enable2FAForSuperUser(c.Request.Context(), id, body.Secret); err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to enable 2FA", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "2FA enabled", nil, nil)
	c.JSON(http.StatusOK, response)
}

// Disable 2FA for SuperUser
func (h *SuperUserGinHandler) Disable2FAForSuperUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.service.Disable2FAForSuperUser(c.Request.Context(), id); err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to disable 2FA", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "2FA disabled", nil, nil)
	c.JSON(http.StatusOK, response)
}

// Get all 2FA-enabled SuperUsers
func (h *SuperUserGinHandler) GetAll2FAEnabledSuperUsersHandler(c *gin.Context) {
	superUsers, err := h.service.GetAll2FAEnabledSuperUsers(c.Request.Context())
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to retrieve SuperUsers", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "SuperUsers retrieved successfully", superUsers, nil)
	c.JSON(http.StatusOK, response)
}

// Update SuperUser role
func (h *SuperUserGinHandler) UpdateSuperUserRoleHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Role == "" {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.service.UpdateSuperUserRole(c.Request.Context(), id, body.Role); err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to update role", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "Role updated", nil, nil)
	c.JSON(http.StatusOK, response)
}

// Update SuperUser permissions
func (h *SuperUserGinHandler) UpdateSuperUserPermissionsHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var body struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.Permissions) == 0 {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.service.UpdateSuperUserPermissions(c.Request.Context(), id, body.Permissions); err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to update permissions", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "Permissions updated", nil, nil)
	c.JSON(http.StatusOK, response)
}

// Update specific SuperUser field
func (h *SuperUserGinHandler) UpdateSuperUserFieldHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	field := c.Param("field")
	var value interface{}
	if err := c.ShouldBindJSON(&value); err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.service.UpdateSuperUserField(c.Request.Context(), id, field, value); err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to update field", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "Field updated", nil, nil)
	c.JSON(http.StatusOK, response)
}

// Generate and set reset token
func (h *SuperUserGinHandler) GenerateAndSetResetTokenHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.service.GenerateAndSetResetToken(c.Request.Context(), id)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to generate reset token", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "Reset token generated", gin.H{"reset_token": token}, nil)
	c.JSON(http.StatusOK, response)
}

// Clear reset token
func (h *SuperUserGinHandler) ClearResetTokenHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.service.ClearResetToken(c.Request.Context(), id); err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to clear reset token", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "Reset token cleared", nil, nil)
	c.JSON(http.StatusOK, response)
}

// Delete SuperUser
func (h *SuperUserGinHandler) DeleteSuperUserByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.service.DeleteSuperUserByID(c.Request.Context(), id); err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to delete SuperUser", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "SuperUser deleted successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}

// Search SuperUsers
func (h *SuperUserGinHandler) SearchSuperUsersHandler(c *gin.Context) {
	searchQuery := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")

	superUsers, err := h.service.SearchSuperUsers(c.Request.Context(), searchQuery, page, limit, sortBy)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to search SuperUsers", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.NewGinResponse(c, http.StatusOK, "SuperUsers retrieved successfully", superUsers, nil)
	c.JSON(http.StatusOK, response)
}
