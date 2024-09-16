package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/services"
	"github.com/lordofthemind/EventifyGo/internals/types"
)

type SuperUserGinHandler struct {
	service services.SuperUserService
}

func NewSuperUserGinHandler(service services.SuperUserService) *SuperUserGinHandler {
	return &SuperUserGinHandler{service: service}
}

// Create SuperUser handler
func (h *SuperUserGinHandler) CreateSuperUserHandler(c *gin.Context) {
	var superUser types.SuperUserType
	if err := c.ShouldBindJSON(&superUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdSuperUser, err := h.service.CreateSuperUser(c.Request.Context(), &superUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdSuperUser)
}

// Get SuperUser by ID
func (h *SuperUserGinHandler) GetSuperUserByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	superUser, err := h.service.GetSuperUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, superUser)
}

// Get SuperUser by email
func (h *SuperUserGinHandler) GetSuperUserByEmailHandler(c *gin.Context) {
	email := c.Param("email")

	superUser, err := h.service.GetSuperUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, superUser)
}

// Get SuperUser by username
func (h *SuperUserGinHandler) GetSuperUserByUsernameHandler(c *gin.Context) {
	username := c.Param("username")

	superUser, err := h.service.GetSuperUserByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, superUser)
}

// Enable 2FA for SuperUser
func (h *SuperUserGinHandler) Enable2FAForSuperUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var body struct {
		Secret string `json:"secret"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Secret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.Enable2FAForSuperUser(c.Request.Context(), id, body.Secret); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled"})
}

// Disable 2FA for SuperUser
func (h *SuperUserGinHandler) Disable2FAForSuperUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.Disable2FAForSuperUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled"})
}

// Get all 2FA-enabled SuperUsers
func (h *SuperUserGinHandler) GetAll2FAEnabledSuperUsersHandler(c *gin.Context) {
	superUsers, err := h.service.GetAll2FAEnabledSuperUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, superUsers)
}

// Update SuperUser role
func (h *SuperUserGinHandler) UpdateSuperUserRoleHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.UpdateSuperUserRole(c.Request.Context(), id, body.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

// Update SuperUser permissions
func (h *SuperUserGinHandler) UpdateSuperUserPermissionsHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var body struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.Permissions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.UpdateSuperUserPermissions(c.Request.Context(), id, body.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions updated"})
}

// Update specific SuperUser field
func (h *SuperUserGinHandler) UpdateSuperUserFieldHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	field := c.Param("field")
	var value interface{}
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.UpdateSuperUserField(c.Request.Context(), id, field, value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Field updated"})
}

// Generate and set reset token
func (h *SuperUserGinHandler) GenerateAndSetResetTokenHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	token, err := h.service.GenerateAndSetResetToken(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reset_token": token})
}

// Clear reset token
func (h *SuperUserGinHandler) ClearResetTokenHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.ClearResetToken(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset token cleared"})
}

// Delete SuperUser
func (h *SuperUserGinHandler) DeleteSuperUserByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteSuperUserByID(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SuperUser deleted"})
}

// Search SuperUsers
func (h *SuperUserGinHandler) SearchSuperUsersHandler(c *gin.Context) {
	searchQuery := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")

	superUsers, err := h.service.SearchSuperUsers(c.Request.Context(), searchQuery, page, limit, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, superUsers)
}
