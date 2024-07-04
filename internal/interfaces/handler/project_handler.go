package handler

import (
	"net/http"
	"strconv"

	"github.com/kouxi08/Eploy/internal/domain"
	"github.com/kouxi08/Eploy/internal/usecase"

	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	Usecase *usecase.ProjectUsecase
}

func NewProjectHandler(u *usecase.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{Usecase: u}
}

func (h *ProjectHandler) GetProjects(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get("userId").(int)
	projects, err := h.Usecase.GetProjects(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"projects": projects})
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var project domain.Project
	if err := c.Bind(&project); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := c.Request().Context()
	userId := c.Get("userId").(int)
	if err := h.Usecase.CreateProject(ctx, project, userId); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"project": project})
}

func (h *ProjectHandler) GetProjectByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	ctx := c.Request().Context()
	userId := c.Get("userId").(int)
	project, err := h.Usecase.GetProjectByID(ctx, id, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	ctx := c.Request().Context()
	userId := c.Get("userId").(int)
	if err := h.Usecase.DeleteProject(ctx, id, userId); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}
