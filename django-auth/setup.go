package django_auth

//	@title			Swagger django-auth API
//	@version		0.1
//	@description	This is django-auth API OPENAPI Documentation.
//	@termsOfService	http://swagger.io/terms/
//  @BasePath  /api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-APP-TOKEN
//	@description				Description for what is this security definition being used

//	@securityDefinitions.apikey Refresh
//	@in							header
//	@name						X-REFRESH-TOKEN
//	@description				Description for what is this security definition being used

import (
	"github.com/bushubdegefu/m-playground/django-auth/controllers"
	"github.com/bushubdegefu/m-playground/logs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Please Note the sequence you mount the middlewares
func SetupRoutes(app *echo.Echo) {
	logOutput, _ := logs.Logfile("django_auth")

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		Output: logOutput,
	}))

	// then authentication middlware
	gapp := app.Group("/api/v1/django_auth")

	// the Otel spanner middleware
	gapp.Use(otelechospanstarter)

	// db session injection
	gapp.Use(dbsessioninjection)
	gapp.GET("/user", controllers.GetUsers).Name = "django_auth_can_view_user"
	gapp.GET("/user/:user_id", controllers.GetUserByID).Name = "django_auth_can_view_user"
	gapp.POST("/user", controllers.PostUser).Name = "django_auth_can_add_user"
	gapp.PATCH("/user/:user_id", controllers.PatchUser).Name = "django_auth_can_change_user"
	gapp.DELETE("/user/:user_id", controllers.DeleteUser).Name = "django_auth_can_delete_user"

	gapp.POST("/userpermission/:permission_id/:user_id", controllers.AddPermissionToUser).Name = "django_auth_can_add_permission"
	gapp.DELETE("/userpermission/:permission_id/:user_id", controllers.DeletePermissionFromUser).Name = "django_auth_can_delete_permission"
	gapp.GET("/userpermission/:user_id", controllers.GetPermissionsOfUsers).Name = "django_auth_can_view_permission"
	gapp.GET("/permissionnoncomplementuser/:user_id", controllers.GetAllPermissionsOfUsers).Name = "django_auth_can_view_permissioncomplement"
	gapp.GET("/permissioncomplementuser/:user_id", controllers.GetPermissionComplementUsers).Name = "django_auth_can_view_permissioncomplement"

	gapp.POST("/usergroup/:group_id/:user_id", controllers.AddGroupToUser).Name = "django_auth_can_add_group"
	gapp.DELETE("/usergroup/:group_id/:user_id", controllers.DeleteGroupFromUser).Name = "django_auth_can_delete_group"
	gapp.GET("/usergroup/:user_id", controllers.GetGroupsOfUsers).Name = "django_auth_can_view_group"
	gapp.GET("/groupnoncomplementuser/:user_id", controllers.GetAllGroupsOfUsers).Name = "django_auth_can_view_groupcomplement"
	gapp.GET("/groupcomplementuser/:user_id", controllers.GetGroupComplementUsers).Name = "django_auth_can_view_groupcomplement"

	gapp.GET("/group", controllers.GetGroups).Name = "django_auth_can_view_group"
	gapp.GET("/group/:group_id", controllers.GetGroupByID).Name = "django_auth_can_view_group"
	gapp.POST("/group", controllers.PostGroup).Name = "django_auth_can_add_group"
	gapp.PATCH("/group/:group_id", controllers.PatchGroup).Name = "django_auth_can_change_group"
	gapp.DELETE("/group/:group_id", controllers.DeleteGroup).Name = "django_auth_can_delete_group"

	gapp.POST("/grouppermission/:permission_id/:group_id", controllers.AddPermissionToGroup).Name = "django_auth_can_add_permission"
	gapp.DELETE("/grouppermission/:permission_id/:group_id", controllers.DeletePermissionFromGroup).Name = "django_auth_can_delete_permission"
	gapp.GET("/grouppermission/:group_id", controllers.GetPermissionsOfGroups).Name = "django_auth_can_view_permission"
	gapp.GET("/permissionnoncomplementgroup/:group_id", controllers.GetAllPermissionsOfGroups).Name = "django_auth_can_view_permissioncomplement"
	gapp.GET("/permissioncomplementgroup/:group_id", controllers.GetPermissionComplementGroups).Name = "django_auth_can_view_permissioncomplement"

	gapp.GET("/permission", controllers.GetPermissions).Name = "django_auth_can_view_permission"
	gapp.GET("/permission/:permission_id", controllers.GetPermissionByID).Name = "django_auth_can_view_permission"
	gapp.POST("/permission", controllers.PostPermission).Name = "django_auth_can_add_permission"
	gapp.PATCH("/permission/:permission_id", controllers.PatchPermission).Name = "django_auth_can_change_permission"
	gapp.DELETE("/permission/:permission_id", controllers.DeletePermission).Name = "django_auth_can_delete_permission"

}
