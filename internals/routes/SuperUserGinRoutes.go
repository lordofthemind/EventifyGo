package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/EventifyGo/internals/handlers"
)

func SetupSuperUserGinRoutes(r *gin.Engine, handler *handlers.SuperUserGinHandler) {
	r.POST("/superusers", handler.CreateSuperUserHandler)
	r.GET("/superusers", handler.GetAllSuperUsersHandler)
	r.GET("/superusers/:id", handler.GetSuperUserByIDHandler)
	r.GET("/superusers/email/:email", handler.GetSuperUserByEmailHandler)
	r.GET("/superusers/username/:username", handler.GetSuperUserByUsernameHandler)
	r.POST("/superusers/:id/enable2fa", handler.Enable2FAForSuperUserHandler)
	r.POST("/superusers/:id/disable2fa", handler.Disable2FAForSuperUserHandler)
	r.GET("/superusers/2fa", handler.GetAll2FAEnabledSuperUsersHandler)
	r.PUT("/superusers/:id/role", handler.UpdateSuperUserRoleHandler)
	r.PUT("/superusers/:id/permissions", handler.UpdateSuperUserPermissionsHandler)
	r.PUT("/superusers/:id/field/:field", handler.UpdateSuperUserFieldHandler)
	r.POST("/superusers/:id/generate-reset-token", handler.GenerateAndSetResetTokenHandler)
	r.POST("/superusers/:id/clear-reset-token", handler.ClearResetTokenHandler)
	r.DELETE("/superusers/:id", handler.DeleteSuperUserByIDHandler)
	r.GET("/superusers/search", handler.SearchSuperUsersHandler)
}
