package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qsrrvln/mountapi/internal/models"
	"gorm.io/gorm"
)

type RouteHandler struct {
	DB *gorm.DB
}

func NewRouteHandler(db *gorm.DB) *RouteHandler {
	return &RouteHandler{DB: db}
}

// GetRoute func
// @Summary Get a route by ID
// @Description Get details of a specific route
// @Tags routes
// @Produce json
// @Param id path string true "Route ID"
// @Success 200 {object} models.Route
// @Failure 404 {object} map[string]string
// @Router /routes/{id} [get]
func (h *RouteHandler) GetRoute(c *gin.Context) {
	id := c.Param("id")
	var route models.Route
	if err := h.DB.First(&route, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}
	c.JSON(http.StatusOK, route)
}

// GetRoutePosts func
// @Summary Get posts for a route
// @Description Get all posts associated with a route, ordered by index
// @Tags routes
// @Produce json
// @Param id path string true "Route ID"
// @Success 200 {array} models.Post
// @Failure 404 {object} map[string]string
// @Router /routes/{id}/posts [get]
func (h *RouteHandler) GetRoutePosts(c *gin.Context) {
	id := c.Param("id")
	var posts []models.Post
	if err := h.DB.Where("route_id = ?", id).Order("order_index asc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// CreateRoute godoc
// @Summary Create a new route
// @Tags routes
// @Accept json
// @Produce json
// @Param route body models.Route true "Route Data"
// @Security ApiKeyAuth
// @Success 201 {object} models.Route
// @Router /routes [post]
func (h *RouteHandler) CreateRoute(c *gin.Context) {
	var route models.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, route)
}

// UpdateRoute godoc
// @Summary Update a route fully
// @Tags routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID"
// @Param route body models.Route true "Route Data"
// @Security ApiKeyAuth
// @Success 200 {object} models.Route
// @Router /routes/{id} [put]
func (h *RouteHandler) UpdateRoute(c *gin.Context) {
	id := c.Param("id")
	var route models.Route
	if err := h.DB.First(&route, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	var input models.Route
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route.Name = input.Name
	route.StartPointName = input.StartPointName
	route.Difficulty = input.Difficulty
	route.LengthM = input.LengthM
	route.ElevationGainM = input.ElevationGainM
	route.IsOfficial = input.IsOfficial
	route.MountainID = input.MountainID // Allow moving route to another mountain

	if err := h.DB.Save(&route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, route)
}

// PatchRoute godoc
// @Summary Update a route partially
// @Tags routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID"
// @Param route body map[string]interface{} true "Fields to update"
// @Security ApiKeyAuth
// @Success 200 {object} models.Route
// @Router /routes/{id} [patch]
func (h *RouteHandler) PatchRoute(c *gin.Context) {
	id := c.Param("id")
	var route models.Route
	if err := h.DB.First(&route, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Model(&route).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, route)
}

// DeleteRoute godoc
// @Summary Delete a route
// @Tags routes
// @Param id path string true "Route ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Router /routes/{id} [delete]
func (h *RouteHandler) DeleteRoute(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Where("id = ?", id).Delete(&models.Route{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
