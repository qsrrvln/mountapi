package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qsrrvln/mountapi/internal/models"
	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{DB: db}
}

// CreatePost godoc
// @Summary Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post Data"
// @Security ApiKeyAuth
// @Success 201 {object} models.Post
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, post)
}

// UpdatePost godoc
// @Summary Update a post fully
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body models.Post true "Post Data"
// @Security ApiKeyAuth
// @Success 200 {object} models.Post
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := h.DB.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var input models.Post
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Name = input.Name
	post.OrderIndex = input.OrderIndex
	post.AltitudeM = input.AltitudeM
	post.DistanceFromStartM = input.DistanceFromStartM
	post.SegmentToNextM = input.SegmentToNextM
	post.DistanceToSummitM = input.DistanceToSummitM
	post.RouteID = input.RouteID

	if err := h.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// PatchPost godoc
// @Summary Update a post partially
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body map[string]interface{} true "Fields to update"
// @Security ApiKeyAuth
// @Success 200 {object} models.Post
// @Router /posts/{id} [patch]
func (h *PostHandler) PatchPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := h.DB.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Model(&post).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// DeletePost godoc
// @Summary Delete a post
// @Tags posts
// @Param id path string true "Post ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Where("id = ?", id).Delete(&models.Post{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// BulkUploadPosts godoc
// @Summary Bulk upload posts from CSV
// @Description Upload multiple posts given a route_id via CSV
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param route_id formData string true "Route ID"
// @Param file formData file true "CSV File (name,order_index,altitude_m,distance_to_summit_m,...)"
// @Security ApiKeyAuth
// @Success 201 {object} map[string]interface{}
// @Router /posts/bulk [post]
func (h *PostHandler) BulkUploadPosts(c *gin.Context) {
	routeIDStr := c.PostForm("route_id")
	if routeIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "route_id is required"})
		return
	}
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid (uuid) route_id"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse csv"})
		return
	}

	var insertedCount int
	// Assuming CSV header: name, order_index, altitude_m, distance_to_summit_m, segment_to_next_m, distance_from_start_m
	for i, record := range records {
		if i == 0 {
			continue
		} // Skip header

		if len(record) < 4 {
			continue
		} // Ensure minimum fields

		orderIdx, _ := strconv.Atoi(record[1])
		alt, _ := strconv.Atoi(record[2])
		distSummit, _ := strconv.Atoi(record[3])

		segNext := 0
		if len(record) > 4 {
			segNext, _ = strconv.Atoi(record[4])
		}

		distStart := 0
		if len(record) > 5 {
			distStart, _ = strconv.Atoi(record[5])
		}

		post := models.Post{
			RouteID:            routeID,
			Name:               record[0],
			OrderIndex:         orderIdx,
			AltitudeM:          alt,
			DistanceToSummitM:  distSummit,
			SegmentToNextM:     segNext,
			DistanceFromStartM: distStart,
		}

		if err := h.DB.Create(&post).Error; err == nil {
			insertedCount++
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bulk upload processed", "inserted": insertedCount})
}
