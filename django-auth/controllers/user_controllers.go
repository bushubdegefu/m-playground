package controllers

import (
	"net/http"
	"strconv"

	"github.com/bushubdegefu/m-playground/common"
	"github.com/bushubdegefu/m-playground/django-auth/models"
	"github.com/bushubdegefu/m-playground/django-auth/services"
	"github.com/bushubdegefu/m-playground/observe"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// GetUsers function to get a Users with pagination and searchFields
// @Summary Get Users
// @Description Get Users
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Param username query string false "Search by username optional field string"
// @Param email query string false "Search by email optional field string"
// @Param first_name query string false "Search by first_name optional field string"
// @Param last_name query string false "Search by last_name optional field string"
// @Success 200 {object} common.ResponsePagination{data=[]models.UserGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /django_auth/user [get]
func GetUsers(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.QueryParam("page"))
	Limit, _ := strconv.Atoi(contx.QueryParam("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	// Getting search fields
	searchTerm := make(map[string]any)
	if err := contx.Bind(&searchTerm); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	searchFields := []string{"username", "email", "first_name", "last_name"}
	filteredSearchTerm := common.FilterSearchTerms(searchTerm, searchFields)

	// Prepare pagination model
	pagination := models.Pagination{
		Page: Page - 1, // assuming pages are 0-indexed in backend
		Size: Limit,
	}

	// Fetch users from service
	users, totalCount, err := services.HandlerUserService.Get(tracer.Tracer, pagination, searchFields, filteredSearchTerm)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// Send paginated response
	return contx.JSON(http.StatusOK, common.ResponsePagination{
		Success: true,
		Message: "Success.",
		Items:   users,
		Total:   totalCount,
		Page:    uint(Page),
		Size:    uint(Limit),
	})
}

// GetUserByID is a function to get a Users by ID
// @Summary Get User by ID
// @Description Get user by ID
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} common.ResponseHTTP{data=models.UserGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /django_auth/user/{user_id} [get]
func GetUserByID(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  parsing Query Prameters
	id := contx.Param("user_id")

	// Fetch user from service
	user, err := services.HandlerUserService.GetOne(tracer.Tracer, id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// Send paginated response
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success",
		Data:    user,
	})

}

// Add User to data
// @Summary Add a new User
// @Description Add User
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body models.UserPost true "Add User"
// @Success 200 {object} common.ResponseHTTP{data=models.UserPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/user [post]
func PostUser(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_user := new(models.UserPost)

	//first parse request data
	if err := contx.Bind(&posted_user); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_user); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// post user from service
	user, err := services.HandlerUserService.Create(tracer.Tracer, posted_user)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "User created successfully.",
		Data:    user,
	})
}

// Patch User to data
// @Summary Patch User
// @Description Patch User
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body models.UserPatch true "Patch User"
// @Param user_id path string true "User ID"
// @Success 200 {object} common.ResponseHTTP{data=models.UserPatch}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/user/{user_id} [patch]
func PatchUser(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validator initialization
	validate := validator.New()

	//getting object_id from path param
	// validate path params
	id := contx.Param("user_id")

	// validate data struct
	patch_user := new(models.UserPatch)
	if err := contx.Bind(&patch_user); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(patch_user); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// patch user from service
	user, err := services.HandlerUserService.Update(tracer.Tracer, patch_user, id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "User updated successfully.",
		Data:    user,
	})
}

// DeleteUsers function removes a user by ID
// @Summary Remove User by ID
// @Description Remove user by ID
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /django_auth/user/{user_id} [delete]
func DeleteUser(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	id := contx.Param("user_id")

	// delete user from service
	err := services.HandlerUserService.Delete(tracer.Tracer, id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// Return success respons
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "User deleted successfully.",
		Data:    nil,
	})
}

// ##########################################################
// ##########  Relationship  Services to Permission
// ##########################################################
// Add Permission to User
// @Summary Add User to Permission
// @Description Add Permission User
// @Tags PermissionUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param permission_id path string true "Permission ID"
// @Param user_id path string true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/userpermission/{permission_id}/{user_id} [post]
func AddPermissionToUser(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	permission_id := contx.Param("permission_id")

	// validate path params
	user_id := contx.Param("user_id")

	err := services.HandlerUserService.AddUserToPermission(tracer.Tracer, user_id, permission_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Added Permission to  User.",
		Data:    nil,
	})
}

// Delete Permission from User
// @Summary Delete Permission
// @Description Delete Permission User
// @Tags PermissionUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param permission_id path string true "Permission ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} common.ResponseHTTP{data=models.UserPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/userpermission/{permission_id}/{user_id} [delete]
func DeletePermissionFromUser(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	permission_id := contx.Param("permission_id")

	// validate path params
	user_id := contx.Param("user_id")

	// removing PermissionFromUser
	err := services.HandlerUserService.RemoveUserFromPermission(tracer.Tracer, user_id, permission_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Removing Permission From User.",
		Data:    nil,
	})
}

// Get Permissions of User
// @Summary Get User to Permission
// @Description Get Permission User
// @Tags PermissionUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.PermissionGet}
// @Param user_id path string true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/userpermission/{user_id} [get]
func GetPermissionsOfUsers(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.QueryParam("page"))
	Limit, _ := strconv.Atoi(contx.QueryParam("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	// validate path params
	user_id := contx.Param("user_id")

	// Prepare pagination model
	pagination := models.Pagination{
		Page: Page - 1, // assuming pages are 0-indexed in backend
		Size: Limit,
	}

	// Fetch users from service
	permissions, totalCount, err := services.HandlerUserService.GetUserPermissions(tracer.Tracer, user_id, pagination)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// Send paginated response
	return contx.JSON(http.StatusOK, common.ResponsePagination{
		Success: true,
		Message: "Success",
		Items:   permissions,
		Total:   totalCount,
		Page:    uint(Page),
		Size:    uint(Limit),
	})
}

// #########################
// No Pagination Services###
// #########################

// Get Permissions of User Complement
// @Summary Get User to Permission Complement
// @Description Get Permission User Complement
// @Tags PermissionUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=[]models.PermissionGet}
// @Param user_id path string true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/permissionnoncomplementuser/{user_id} [get]
func GetAllPermissionsOfUsers(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	user_id := contx.Param("user_id")

	// Fetch users from service
	permissions, err := services.HandlerUserService.GetAllPermissionsForUser(tracer.Tracer, user_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Data:    permissions,
		Message: "working",
	})
}

// Get Permissions of User Not Complement
// @Summary Get User to Permission Not Complement
// @Description Get Permission User Not Complement
// @Tags PermissionUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=[]models.PermissionGet}
// @Param user_id path string true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/permissioncomplementuser/{user_id} [get]
func GetPermissionComplementUsers(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	user_id := contx.Param("user_id")

	// Fetch users from service
	permissions, err := services.HandlerUserService.GetAllPermissionsuserDoesNotHave(tracer.Tracer, user_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Data:    permissions,
		Message: "working",
	})
}

// ##########################################################
// ##########  Relationship  Services to Group
// ##########################################################
// Add Group to User
// @Summary Add User to Group
// @Description Add Group User
// @Tags GroupUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Param user_id path string true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/usergroup/{group_id}/{user_id} [post]
func AddGroupToUser(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	group_id := contx.Param("group_id")

	// validate path params
	user_id := contx.Param("user_id")

	err := services.HandlerUserService.AddUserToGroup(tracer.Tracer, user_id, group_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Added Group to  User.",
		Data:    nil,
	})
}

// Delete Group from User
// @Summary Delete Group
// @Description Delete Group User
// @Tags GroupUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} common.ResponseHTTP{data=models.UserPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/usergroup/{group_id}/{user_id} [delete]
func DeleteGroupFromUser(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	group_id := contx.Param("group_id")

	// validate path params
	user_id := contx.Param("user_id")

	// removing GroupFromUser
	err := services.HandlerUserService.RemoveUserFromGroup(tracer.Tracer, user_id, group_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Removing Group From User.",
		Data:    nil,
	})
}

// Get Groups of User
// @Summary Get User to Group
// @Description Get Group User
// @Tags GroupUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.GroupGet}
// @Param user_id path string true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/usergroup/{user_id} [get]
func GetGroupsOfUsers(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.QueryParam("page"))
	Limit, _ := strconv.Atoi(contx.QueryParam("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	// validate path params
	user_id := contx.Param("user_id")

	// Prepare pagination model
	pagination := models.Pagination{
		Page: Page - 1, // assuming pages are 0-indexed in backend
		Size: Limit,
	}

	// Fetch users from service
	groups, totalCount, err := services.HandlerUserService.GetUserGroups(tracer.Tracer, user_id, pagination)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// Send paginated response
	return contx.JSON(http.StatusOK, common.ResponsePagination{
		Success: true,
		Message: "Success",
		Items:   groups,
		Total:   totalCount,
		Page:    uint(Page),
		Size:    uint(Limit),
	})
}

// #########################
// No Pagination Services###
// #########################

// Get Groups of User Complement
// @Summary Get User to Group Complement
// @Description Get Group User Complement
// @Tags GroupUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=[]models.GroupGet}
// @Param user_id path string true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/groupnoncomplementuser/{user_id} [get]
func GetAllGroupsOfUsers(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	user_id := contx.Param("user_id")

	// Fetch users from service
	groups, err := services.HandlerUserService.GetAllGroupsForUser(tracer.Tracer, user_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Data:    groups,
		Message: "working",
	})
}

// Get Groups of User Not Complement
// @Summary Get User to Group Not Complement
// @Description Get Group User Not Complement
// @Tags GroupUsers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=[]models.GroupGet}
// @Param user_id path string true "User ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/groupcomplementuser/{user_id} [get]
func GetGroupComplementUsers(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	user_id := contx.Param("user_id")

	// Fetch users from service
	groups, err := services.HandlerUserService.GetAllGroupsuserDoesNotHave(tracer.Tracer, user_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Data:    groups,
		Message: "working",
	})
}
