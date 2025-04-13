package handler

import (
	"strconv"

	"github.com/Mukam21/server_Golang/pkg/model"
	"github.com/Mukam21/server_Golang/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	service *service.Service
	log     *logrus.Logger
}

func NewHandler(service *service.Service, log *logrus.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) InitRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		persons := api.Group("/persons")
		{
			persons.POST("", h.createPerson)
			persons.GET("", h.getPersons)
			persons.GET("/:id", h.getPerson)
			persons.PUT("/:id", h.updatePerson)
			persons.PATCH("/:id", h.patchPerson)
			persons.DELETE("/:id", h.deletePerson)
		}
	}
}

// @Summary Create a new person
// @Description Create a person with name, surname, and optional patronymic
// @Tags persons
// @Accept json
// @Produce json
// @Param person body model.PersonRequest true "Person data"
// @Success 201 {object} model.Person
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/persons [post]
func (h *Handler) createPerson(c *gin.Context) {
	var req model.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug("Invalid request: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	person, err := h.service.CreatePerson(&req)
	if err != nil {
		h.log.Errorf("Failed to create person: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, person)
}

// @Summary Get list of persons
// @Description Retrieve persons with pagination and optional filters
// @Tags persons
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param name query string false "Filter by name"
// @Param surname query string false "Filter by surname"
// @Param age query int false "Filter by age"
// @Param gender query string false "Filter by gender" enum(male,female,other)
// @Param nationality query string false "Filter by nationality"
// @Success 200 {array} model.Person
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/persons [get]
func (h *Handler) getPersons(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	filters := map[string]string{
		"name":        c.Query("name"),
		"surname":     c.Query("surname"),
		"age":         c.Query("age"),
		"gender":      c.Query("gender"),
		"nationality": c.Query("nationality"),
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.log.Debug("Invalid page number: ", pageStr)
		c.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		h.log.Debug("Invalid limit: ", limitStr)
		c.JSON(400, gin.H{"error": "Invalid limit"})
		return
	}

	persons, err := h.service.GetAll(page, limit, filters)
	if err != nil {
		h.log.Errorf("Failed to get persons: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, persons)
}

// @Summary Get person by ID
// @Description Retrieve a person by their ID
// @Tags persons
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} model.Person
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/persons/{id} [get]
func (h *Handler) getPerson(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Debug("Invalid ID: ", c.Param("id"))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	person, err := h.service.GetByID(id)
	if err != nil {
		h.log.Errorf("Failed to get person with ID %d: %v", id, err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if person == nil {
		c.JSON(404, gin.H{"error": "Person not found"})
		return
	}
	c.JSON(200, person)
}

// @Summary Update a person
// @Description Update person details by ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body model.Person true "Updated person data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/persons/{id} [put]
func (h *Handler) updatePerson(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Debug("Invalid ID: ", c.Param("id"))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var person model.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		h.log.Debug("Invalid request: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	person.ID = id

	if err := h.service.Update(&person); err != nil {
		h.log.Errorf("Failed to update person with ID %d: %v", id, err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Person updated"})
}

// @Summary Partially update a person
// @Description Update specific fields of a person by ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body model.PersonPatchRequest true "Fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/persons/{id} [patch]
func (h *Handler) patchPerson(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Debug("Invalid ID: ", c.Param("id"))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var patch model.PersonPatchRequest
	if err := c.ShouldBindJSON(&patch); err != nil {
		h.log.Debug("Invalid request: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Patch(id, &patch); err != nil {
		h.log.Errorf("Failed to patch person with ID %d: %v", id, err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Person patched"})
}

// @Summary Delete a person
// @Description Delete a person by ID
// @Tags persons
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/persons/{id} [delete]
func (h *Handler) deletePerson(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Debug("Invalid ID: ", c.Param("id"))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		h.log.Errorf("Failed to delete person with ID %d: %v", id, err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Person deleted"})
}
