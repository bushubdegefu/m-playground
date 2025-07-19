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

// GetGroups function to get a Groups with pagination and searchFields
// @Summary Get Groups
// @Description Get Groups
// @Tags Groups
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Param name query string false "Search by name optional field string"
// @Success 200 {object} common.ResponsePagination{data=[]models.GroupGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /django_auth/group [get]
func GetGroups(contx echo.Context) error {
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

	// Fetch groups from service
	groups, totalCount, err := services.HandlerGroupService.Get(tracer.Tracer, pagination, searchFields, filteredSearchTerm)
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
		Items:   groups,
		Total:   totalCount,
		Page:    uint(Page),
		Size:    uint(Limit),
	})
}

// GetGroupByID is a function to get a Groups by ID
// @Summary Get Group by ID
// @Description Get group by ID
// @Tags Groups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Success 200 {object} common.ResponseHTTP{data=models.GroupGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /django_auth/group/{group_id} [get]
func GetGroupByID(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  parsing Query Prameters
	id := contx.Param("group_id")

	// Fetch group from service
	group, err := services.HandlerGroupService.GetOne(tracer.Tracer, id)
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
		Data:    group,
	})

}

// Add Group to data
// @Summary Add a new Group
// @Description Add Group
// @Tags Groups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param group body models.GroupPost true "Add Group"
// @Success 200 {object} common.ResponseHTTP{data=models.GroupPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/group [post]
func PostGroup(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_group := new(models.GroupPost)

	//first parse request data
	if err := contx.Bind(&posted_group); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_group); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// post group from service
	group, err := services.HandlerGroupService.Create(tracer.Tracer, posted_group)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Group created successfully.",
		Data:    group,
	})
}

// Patch Group to data
// @Summary Patch Group
// @Description Patch Group
// @Tags Groups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param group body models.GroupPatch true "Patch Group"
// @Param group_id path string true "Group ID"
// @Success 200 {object} common.ResponseHTTP{data=models.GroupPatch}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/group/{group_id} [patch]
func PatchGroup(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validator initialization
	validate := validator.New()

	//getting object_id from path param
	// validate path params
	id := contx.Param("group_id")

	// validate data struct
	patch_group := new(models.GroupPatch)
	if err := contx.Bind(&patch_group); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(patch_group); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// patch group from service
	group, err := services.HandlerGroupService.Update(tracer.Tracer, patch_group, id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Group updated successfully.",
		Data:    group,
	})
}

// DeleteGroups function removes a group by ID
// @Summary Remove Group by ID
// @Description Remove group by ID
// @Tags Groups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /django_auth/group/{group_id} [delete]
func DeleteGroup(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	id := contx.Param("group_id")

	// delete group from service
	err := services.HandlerGroupService.Delete(tracer.Tracer, id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// Return success respons
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Group deleted successfully.",
		Data:    nil,
	})
}

// ##########################################################
// ##########  Relationship  Services to Permission
// ##########################################################
// Add Permission to Group
// @Summary Add Group to Permission
// @Description Add Permission Group
// @Tags PermissionGroups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param permission_id path string true "Permission ID"
// @Param group_id path string true "Group ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/grouppermission/{permission_id}/{group_id} [post]
func AddPermissionToGroup(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	permission_id := contx.Param("permission_id")

	// validate path params
	group_id := contx.Param("group_id")

	err := services.HandlerGroupService.AddGroupToPermission(tracer.Tracer, group_id, permission_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Added Permission to  Group.",
		Data:    nil,
	})
}

// Delete Permission from Group
// @Summary Delete Permission
// @Description Delete Permission Group
// @Tags PermissionGroups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param permission_id path string true "Permission ID"
// @Param group_id path string true "Group ID"
// @Success 200 {object} common.ResponseHTTP{data=models.GroupPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/grouppermission/{permission_id}/{group_id} [delete]
func DeletePermissionFromGroup(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	permission_id := contx.Param("permission_id")

	// validate path params
	group_id := contx.Param("group_id")

	// removing PermissionFromGroup
	err := services.HandlerGroupService.RemoveGroupFromPermission(tracer.Tracer, group_id, permission_id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Removing Permission From Group.",
		Data:    nil,
	})
}

// Get Permissions of Group
// @Summary Get Group to Permission
// @Description Get Permission Group
// @Tags PermissionGroups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.PermissionGet}
// @Param group_id path string true "Group ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/grouppermission/{group_id} [get]
func GetPermissionsOfGroups(contx echo.Context) error {
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
	group_id := contx.Param("group_id")

	// Prepare pagination model
	pagination := models.Pagination{
		Page: Page - 1, // assuming pages are 0-indexed in backend
		Size: Limit,
	}

	// Fetch groups from service
	permissions, totalCount, err := services.HandlerGroupService.GetGroupPermissions(tracer.Tracer, group_id, pagination)
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

// Get Permissions of Group Complement
// @Summary Get Group to Permission Complement
// @Description Get Permission Group Complement
// @Tags PermissionGroups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=[]models.PermissionGet}
// @Param group_id path string true "Group ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/permissionnoncomplementgroup/{group_id} [get]
func GetAllPermissionsOfGroups(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	group_id := contx.Param("group_id")

	// Fetch groups from service
	permissions, err := services.HandlerGroupService.GetAllPermissionsForGroup(tracer.Tracer, group_id)
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

// Get Permissions of Group Not Complement
// @Summary Get Group to Permission Not Complement
// @Description Get Permission Group Not Complement
// @Tags PermissionGroups
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=[]models.PermissionGet}
// @Param group_id path string true "Group ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /django_auth/permissioncomplementgroup/{group_id} [get]
func GetPermissionComplementGroups(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	group_id := contx.Param("group_id")

	// Fetch groups from service
	permissions, err := services.HandlerGroupService.GetAllPermissionsgroupDoesNotHave(tracer.Tracer, group_id)
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
