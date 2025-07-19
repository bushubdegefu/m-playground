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

// GetPermissions function to get a Permissions with pagination and searchFields
// @Summary Get Permissions
// @Description Get Permissions
// @Tags Permissions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Param name query string false "Search by name optional field string"
// @Param codename query string false "Search by codename optional field string"
// @Success 200 {object} common.ResponsePagination{data=[]models.PermissionGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /django_auth/permission [get]
func GetPermissions(contx echo.Context) error {
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

	// Fetch permissions from service
	permissions, totalCount, err := services.HandlerPermissionService.Get(tracer.Tracer, pagination, searchFields, filteredSearchTerm)
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
		Items:   permissions,
		Total:   totalCount,
		Page:    uint(Page),
		Size:    uint(Limit),
	})
}

// GetPermissionByID is a function to get a Permissions by ID
// @Summary Get Permission by ID
// @Description Get permission by ID
// @Tags Permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} common.ResponseHTTP{data=models.PermissionGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /django_auth/permission/{permission_id} [get]
func GetPermissionByID(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  parsing Query Prameters
	id := contx.Param("permission_id")

	// Fetch permission from service
	permission, err := services.HandlerPermissionService.GetOne(tracer.Tracer, id)
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
		Data:    permission,
	})

}

// Add Permission to data
// @Summary Add a new Permission
// @Description Add Permission
// @Tags Permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param permission body models.PermissionPost true "Add Permission"
// @Success 200 {object} common.ResponseHTTP{data=models.PermissionPost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/permission [post]
func PostPermission(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_permission := new(models.PermissionPost)

	//first parse request data
	if err := contx.Bind(&posted_permission); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_permission); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// post permission from service
	permission, err := services.HandlerPermissionService.Create(tracer.Tracer, posted_permission)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Permission created successfully.",
		Data:    permission,
	})
}

// Patch Permission to data
// @Summary Patch Permission
// @Description Patch Permission
// @Tags Permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param permission body models.PermissionPatch true "Patch Permission"
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} common.ResponseHTTP{data=models.PermissionPatch}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /django_auth/permission/{permission_id} [patch]
func PatchPermission(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validator initialization
	validate := validator.New()

	//getting object_id from path param
	// validate path params
	id := contx.Param("permission_id")

	// validate data struct
	patch_permission := new(models.PermissionPatch)
	if err := contx.Bind(&patch_permission); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(patch_permission); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// patch permission from service
	permission, err := services.HandlerPermissionService.Update(tracer.Tracer, patch_permission, id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Permission updated successfully.",
		Data:    permission,
	})
}

// DeletePermissions function removes a permission by ID
// @Summary Remove Permission by ID
// @Description Remove permission by ID
// @Tags Permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /django_auth/permission/{permission_id} [delete]
func DeletePermission(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	// validate path params
	id := contx.Param("permission_id")

	// delete permission from service
	err := services.HandlerPermissionService.Delete(tracer.Tracer, id)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
		})
	}

	// Return success respons
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Permission deleted successfully.",
		Data:    nil,
	})
}
