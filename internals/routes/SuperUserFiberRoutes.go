package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordofthemind/EventifyGo/internals/handlers"
)

func SetupSuperUserFiberRoutes(app *fiber.App, handler *handlers.SuperUserFiberHandler) {
	app.Post("/superusers", handler.CreateSuperUserHandler)
	app.Get("/superusers", handler.GetAllSuperUsersHandler)
	app.Get("/superusers/:id", handler.GetSuperUserByIDHandler)
	app.Get("/superusers/email/:email", handler.GetSuperUserByEmailHandler)
	app.Get("/superusers/username/:username", handler.GetSuperUserByUsernameHandler)
	app.Post("/superusers/:id/enable2fa", handler.Enable2FAForSuperUserHandler)
	app.Post("/superusers/:id/disable2fa", handler.Disable2FAForSuperUserHandler)
	app.Get("/superusers/2fa", handler.GetAll2FAEnabledSuperUsersHandler)
	app.Put("/superusers/:id/role", handler.UpdateSuperUserRoleHandler)
	app.Put("/superusers/:id/permissions", handler.UpdateSuperUserPermissionsHandler)
	app.Put("/superusers/:id/field/:field", handler.UpdateSuperUserFieldHandler)
	app.Post("/superusers/:id/generate-reset-token", handler.GenerateAndSetResetTokenHandler)
	app.Post("/superusers/:id/clear-reset-token", handler.ClearResetTokenHandler)
	app.Delete("/superusers/:id", handler.DeleteSuperUserByIDHandler)
	app.Get("/superusers/search", handler.SearchSuperUsersHandler)
}
