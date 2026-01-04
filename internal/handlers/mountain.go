package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qsrrvln/mountapi/internal/models"
	"gorm.io/gorm"
)

type MountainHandler struct {
	DB *gorm.DB
}

func NewMountainHandler(db *gorm.DB) *MountainHandler {
	return &MountainHandler{DB: db}
}

// GetMountains godoc
// @Summary Get all mountains
// @Description Get a list of mountains with optional filtering
// @Tags mountains
// @Produce json
// @Param query query string false "Search query"
// @Param province query string false "Filter by province"
// @Success 200 {array} models.Mountain
// @Router /mountains [get]
func (h *MountainHandler) GetMountains(c *gin.Context) {
	var mountains []models.Mountain
	query := h.DB.Model(&models.Mountain{})

	if q := c.Query("query"); q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}
	if province := c.Query("province"); province != "" {
		query = query.Where("province = ?", province)
	}

	if err := query.Find(&mountains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mountains)
}

// GetMountain func
// @Summary Get a mountain by ID
// @Description Get details of a specific mountain
// @Tags mountains
// @Produce json
// @Param id path string true "Mountain ID"
// @Success 200 {object} models.Mountain
// @Failure 404 {object} map[string]string
// @Router /mountains/{id} [get]
func (h *MountainHandler) GetMountain(c *gin.Context) {
	id := c.Param("id")
	var mountain models.Mountain
	if err := h.DB.Preload("Routes").First(&mountain, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mountain not found"})
		return
	}
	c.JSON(http.StatusOK, mountain)
}

// GetMountainRoutes func
// @Summary Get routes for a mountain
// @Description Get all routes associated with a mountain
// @Tags mountains
// @Produce json
// @Param id path string true "Mountain ID"
// @Success 200 {array} models.Route
// @Failure 404 {object} map[string]string
// @Router /mountains/{id}/routes [get]
func (h *MountainHandler) GetMountainRoutes(c *gin.Context) {
	id := c.Param("id")
	var routes []models.Route
	if err := h.DB.Where("mountain_id = ?", id).Find(&routes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, routes)
}

// CreateMountain godoc
// @Summary Create a new mountain
// @Description Create a new mountain record
// @Tags mountains
// @Accept json
// @Produce json
// @Param mountain body models.Mountain true "Mountain Data"
// @Security ApiKeyAuth
// @Success 201 {object} models.Mountain
// @Failure 400 {object} map[string]string
// @Router /mountains [post]
func (h *MountainHandler) CreateMountain(c *gin.Context) {
	var mountain models.Mountain
	if err := c.ShouldBindJSON(&mountain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&mountain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mountain)
}

// UpdateMountain godoc
// @Summary Update a mountain fully
// @Description Update all fields of a mountain
// @Tags mountains
// @Accept json
// @Produce json
// @Param id path string true "Mountain ID"
// @Param mountain body models.Mountain true "Mountain Data"
// @Security ApiKeyAuth
// @Success 200 {object} models.Mountain
// @Failure 404 {object} map[string]string
// @Router /mountains/{id} [put]
func (h *MountainHandler) UpdateMountain(c *gin.Context) {
	id := c.Param("id")
	var mountain models.Mountain
	if err := h.DB.First(&mountain, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mountain not found"})
		return
	}

	var input models.Mountain
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Full update
	mountain.Name = input.Name
	mountain.ElevationM = input.ElevationM
	mountain.Province = input.Province
	mountain.Latitude = input.Latitude
	mountain.Longitude = input.Longitude

	if err := h.DB.Save(&mountain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mountain)
}

// PatchMountain godoc
// @Summary Update a mountain partially
// @Description Update specific fields of a mountain
// @Tags mountains
// @Accept json
// @Produce json
// @Param id path string true "Mountain ID"
// @Param mountain body map[string]interface{} true "Fields to update"
// @Security ApiKeyAuth
// @Success 200 {object} models.Mountain
// @Failure 404 {object} map[string]string
// @Router /mountains/{id} [patch]
func (h *MountainHandler) PatchMountain(c *gin.Context) {
	id := c.Param("id")
	var mountain models.Mountain
	if err := h.DB.First(&mountain, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mountain not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Model(&mountain).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mountain)
}

// DeleteMountain godoc
// @Summary Delete a mountain
// @Description Delete a mountain and its routes/posts
// @Tags mountains
// @Param id path string true "Mountain ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /mountains/{id} [delete]
func (h *MountainHandler) DeleteMountain(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Where("id = ?", id).Delete(&models.Mountain{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
