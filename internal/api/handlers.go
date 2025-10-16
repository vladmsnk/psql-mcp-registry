package api

import (
	"errors"
	"net/http"

	"psql-mcp-registry/internal/instance_manager"
	"psql-mcp-registry/internal/model"

	"github.com/gin-gonic/gin"
)

// RegisterInstance handles POST /api/v1/instances
func (s *APIServer) RegisterInstance(c *gin.Context) {
	var req RegisterInstanceRequest

	// Parse and validate JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Create instance model with default status
	instance := model.Instance{
		Name:            req.Name,
		DatabaseName:    req.DatabaseName,
		Description:     req.Description,
		CreatorUsername: req.CreatorUsername,
		Status:          "active",
	}

	// Register the instance
	err := s.manager.RegisterInstance(c.Request.Context(), instance)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, instance_manager.ErrInstanceAlreadyExists) {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "instance_already_exists",
				Message: "An instance with this name already exists",
			})
			return
		}

		// Handle other errors
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
		return
	}

	// Get the created instance to return full details
	createdInstance, err := s.manager.GetInstance(c.Request.Context(), instance.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "failed_to_retrieve_instance",
			Message: err.Error(),
		})
		return
	}

	// Return success response
	response := RegisterInstanceResponse{
		ID:              createdInstance.ID,
		Name:            createdInstance.Name,
		DatabaseName:    createdInstance.DatabaseName,
		Description:     createdInstance.Description,
		CreatorUsername: createdInstance.CreatorUsername,
		Status:          createdInstance.Status,
		CreatedAt:       createdInstance.CreatedAt,
		UpdatedAt:       createdInstance.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// HealthCheck handles GET /health
func (s *APIServer) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

