package handlers

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/services"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"github.com/revel/revel"
)

type SuperUserRevelController struct {
	*revel.Controller
	service services.SuperUserService
}

func (c *SuperUserRevelController) Init(service services.SuperUserService) *SuperUserRevelController {
	return &SuperUserRevelController{service: service}
}

// Create SuperUser handler
func (c *SuperUserRevelController) CreateSuperUserHandler() revel.Result {
	var superUser types.SuperUserType
	if err := c.Params.BindJSON(&superUser); err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid input"})
	}

	createdSuperUser, err := c.service.CreateSuperUser(c.Request.Context(), &superUser)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(createdSuperUser)
}

// Get SuperUser by ID
func (c *SuperUserRevelController) GetSuperUserByIDHandler(id string) revel.Result {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid ID format"})
	}

	superUser, err := c.service.GetSuperUserByID(c.Request.Context(), uuidID)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(superUser)
}

// Get SuperUser by email
func (c *SuperUserRevelController) GetSuperUserByEmailHandler(email string) revel.Result {
	superUser, err := c.service.GetSuperUserByEmail(c.Request.Context(), email)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(superUser)
}

// Get SuperUser by username
func (c *SuperUserRevelController) GetSuperUserByUsernameHandler(username string) revel.Result {
	superUser, err := c.service.GetSuperUserByUsername(c.Request.Context(), username)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(superUser)
}

// Enable 2FA for SuperUser
func (c *SuperUserRevelController) Enable2FAForSuperUserHandler(id string) revel.Result {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid ID format"})
	}

	var body struct {
		Secret string `json:"secret"`
	}
	if err := c.Params.BindJSON(&body); err != nil || body.Secret == "" {
		return c.RenderJSON(map[string]string{"error": "Invalid input"})
	}

	err = c.service.Enable2FAForSuperUser(c.Request.Context(), uuidID, body.Secret)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(map[string]string{"message": "2FA enabled"})
}

// Disable 2FA for SuperUser
func (c *SuperUserRevelController) Disable2FAForSuperUserHandler(id string) revel.Result {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid ID format"})
	}

	err = c.service.Disable2FAForSuperUser(c.Request.Context(), uuidID)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(map[string]string{"message": "2FA disabled"})
}

// Get all 2FA-enabled SuperUsers
func (c *SuperUserRevelController) GetAll2FAEnabledSuperUsersHandler() revel.Result {
	superUsers, err := c.service.GetAll2FAEnabledSuperUsers(c.Request.Context())
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(superUsers)
}

// Update SuperUser role
func (c *SuperUserRevelController) UpdateSuperUserRoleHandler(id string) revel.Result {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid ID format"})
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := c.Params.BindJSON(&body); err != nil || body.Role == "" {
		return c.RenderJSON(map[string]string{"error": "Invalid input"})
	}

	err = c.service.UpdateSuperUserRole(c.Request.Context(), uuidID, body.Role)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(map[string]string{"message": "Role updated"})
}

// Update SuperUser permissions
func (c *SuperUserRevelController) UpdateSuperUserPermissionsHandler(id string) revel.Result {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid ID format"})
	}

	var body struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.Params.BindJSON(&body); err != nil || len(body.Permissions) == 0 {
		return c.RenderJSON(map[string]string{"error": "Invalid input"})
	}

	err = c.service.UpdateSuperUserPermissions(c.Request.Context(), uuidID, body.Permissions)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(map[string]string{"message": "Permissions updated"})
}

// Generate and set reset token
func (c *SuperUserRevelController) GenerateAndSetResetTokenHandler(id string) revel.Result {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid ID format"})
	}

	token, err := c.service.GenerateAndSetResetToken(c.Request.Context(), uuidID)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(map[string]string{"reset_token": token})
}

// Clear reset token
func (c *SuperUserRevelController) ClearResetTokenHandler(id string) revel.Result {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid ID format"})
	}

	err = c.service.ClearResetToken(c.Request.Context(), uuidID)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(map[string]string{"message": "Reset token cleared"})
}

// Delete SuperUser
func (c *SuperUserRevelController) DeleteSuperUserByIDHandler(id string) revel.Result {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": "Invalid ID format"})
	}

	err = c.service.DeleteSuperUserByID(c.Request.Context(), uuidID)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(map[string]string{"message": "SuperUser deleted"})
}

// Search SuperUsers
// Search SuperUsers
func (c *SuperUserRevelController) SearchSuperUsersHandler() revel.Result {
	searchQuery := c.Params.Get("q")

	// Check if "page" parameter exists, otherwise use default value "1"
	pageParam := c.Params.Get("page")
	page := 1
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	// Check if "limit" parameter exists, otherwise use default value "10"
	limitParam := c.Params.Get("limit")
	limit := 10
	if limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}

	// Check if "sortBy" parameter exists, otherwise use default value "created_at"
	sortBy := c.Params.Get("sortBy")
	if sortBy == "" {
		sortBy = "created_at"
	}

	// Call the service with the parsed parameters
	superUsers, err := c.service.SearchSuperUsers(c.Request.Context(), searchQuery, page, limit, sortBy)
	if err != nil {
		return c.RenderJSON(map[string]string{"error": err.Error()})
	}

	return c.RenderJSON(superUsers)
}
