package routes

// import (
// 	"github.com/revel/revel"
// )

// func InitializeRoutes() {
// 	revel.Router.Routes = []*revel.Route{
// 		// Create SuperUser
// 		{
// 			Method:     "POST",
// 			Path:       "/superusers",
// 			Controller: "SuperUserRevelController.CreateSuperUserHandler",
// 		},
// 		// Get SuperUser by ID
// 		{
// 			Method:     "GET",
// 			Path:       "/superusers/:id",
// 			Controller: "SuperUserRevelController.GetSuperUserByIDHandler",
// 		},
// 		// Get SuperUser by email
// 		{
// 			Method:     "GET",
// 			Path:       "/superusers/email/:email",
// 			Controller: "SuperUserRevelController.GetSuperUserByEmailHandler",
// 		},
// 		// Get SuperUser by username
// 		{
// 			Method:     "GET",
// 			Path:       "/superusers/username/:username",
// 			Controller: "SuperUserRevelController.GetSuperUserByUsernameHandler",
// 		},
// 		// Enable 2FA
// 		{
// 			Method:     "POST",
// 			Path:       "/superusers/:id/enable2fa",
// 			Controller: "SuperUserRevelController.Enable2FAForSuperUserHandler",
// 		},
// 		// Disable 2FA
// 		{
// 			Method:     "POST",
// 			Path:       "/superusers/:id/disable2fa",
// 			Controller: "SuperUserRevelController.Disable2FAForSuperUserHandler",
// 		},
// 		// Get all 2FA enabled SuperUsers
// 		{
// 			Method:     "GET",
// 			Path:       "/superusers/2fa",
// 			Controller: "SuperUserRevelController.GetAll2FAEnabledSuperUsersHandler",
// 		},
// 		// Update SuperUser role
// 		{
// 			Method:     "PUT",
// 			Path:       "/superusers/:id/role",
// 			Controller: "SuperUserRevelController.UpdateSuperUserRoleHandler",
// 		},
// 		// Update SuperUser permissions
// 		{
// 			Method:     "PUT",
// 			Path:       "/superusers/:id/permissions",
// 			Controller: "SuperUserRevelController.UpdateSuperUserPermissionsHandler",
// 		},
// 		// Generate and set reset token
// 		{
// 			Method:     "POST",
// 			Path:       "/superusers/:id/generate-reset-token",
// 			Controller: "SuperUserRevelController.GenerateAndSetResetTokenHandler",
// 		},
// 		// Clear reset token
// 		{
// 			Method:     "POST",
// 			Path:       "/superusers/:id/clear-reset-token",
// 			Controller: "SuperUserRevelController.ClearResetTokenHandler",
// 		},
// 		// Delete SuperUser
// 		{
// 			Method:     "DELETE",
// 			Path:       "/superusers/:id",
// 			Controller: "SuperUserRevelController.DeleteSuperUserByIDHandler",
// 		},
// 		// Search SuperUsers
// 		{
// 			Method:     "GET",
// 			Path:       "/superusers/search",
// 			Controller: "SuperUserRevelController.SearchSuperUsersHandler",
// 		},
// 	}
// }
