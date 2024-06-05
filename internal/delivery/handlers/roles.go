package handlers

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/errs"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/services"
	"BACKEND/internal/validators"
	"BACKEND/pkg/responses"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type RoleHandler struct {
	service   services.Base
	converter converters.BaseConverter
	validate  *validator.Validate
}

func InitRoleHandler(
	service services.Base,
	validate *validator.Validate,
) *RoleHandler {
	return &RoleHandler{
		service:   service,
		converter: converters.InitBaseConverter(),
		validate:  validate,
	}
}

// CreateRole
// @Summary Create Role
// @Description Create a new role
// @Tags Roles
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param role body dto.BaseBase true "Role data to create"
// @Success 201 {object} responses.CreatedIDResponse "Role successfully created"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/role [post]
func (r RoleHandler) CreateRole(c *gin.Context) {
	var role dto.BaseBase

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	if err := r.validate.Struct(role); err != nil {
		customErr := validators.CustomErrorMessage(err, &dto.BaseBase{})
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: customErr})
		return
	}

	ctx := c.Request.Context()

	roleID, err := r.service.Create(ctx, r.converter.BaseBaseDTOToDomain(role))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrAlreadyExist):
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.CreatedIDResponse{ID: roleID})
}

// GetRoles
// @Summary Get Roles
// @Description Get roles
// @Tags Roles
// @Accept json
// @Produce json
// @Success 200 {object} []dto.Base "Return roles"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/role [get]
func (r RoleHandler) GetRoles(c *gin.Context) {
	ctx := c.Request.Context()

	roles, err := r.service.GetByName(ctx)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, roles)
}

// DeleteRoles
// @Summary Delete Roles
// @Description Delete roles by provided IDs
// @Tags Roles
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param role_ids body []int true "Role IDs to delete"
// @Success 200 "Roles successfully deleted"
// @Failure 400 {object} responses.MessageResponse "Invalid body or jwt provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/role [delete]
func (r RoleHandler) DeleteRoles(c *gin.Context) {
	var roleIDs []int

	if err := c.ShouldBindJSON(&roleIDs); err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadBody})
		return
	}

	ctx := c.Request.Context()

	err := r.service.Delete(ctx, roleIDs)
	if err != nil {
		if err.Error() == "Один из переданных объектов не существует" {
			c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: err.Error()})
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}
