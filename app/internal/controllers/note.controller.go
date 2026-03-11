package controllers

import (
	"JWT/internal/initializers"
	"JWT/internal/models"
	"JWT/internal/utils"
	"strconv"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/note_v1"
	"github.com/gofiber/fiber/v2"
)

// CreateNote godoc
// @Summary Create a new note
// @Description Creates a new note for the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Param note body desc.CreateNoteRequest true "Note creation request"
// @Success 201 {object} map[string]interface{} "note"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /api/note [post]
func CreateNote(c *fiber.Ctx) error {
	// get id
	targetId, ok := c.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(c)
	}

	// parse body
	request := desc.CreateNoteRequest{}
	if err := utils.ParseBodyAndValidate[desc.CreateNoteRequest](c, &request); err != nil {
		return utils.BadRequest(c)
	}
	request.UserID = targetId

	// create note
	note, err := initializers.GrpcNoteService.CreateNote(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"note": note})
}

// GetAllNotes godoc
// @Summary Get all notes
// @Description Retrieves all notes for the authenticated user with pagination
// @Tags notes
// @Produce json
// @Param limit query int false "Limit number of notes"
// @Param offset query int false "Offset for pagination"
// @Param ascendingOrder query bool false "Sort in ascending order"
// @Success 200 {object} map[string]interface{} "notes"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /api/note [get]
func GetAllNotes(c *fiber.Ctx) error {
	request := models.GetAllItemsRequest{}

	// parse query
	if err := c.QueryParser(&request); err != nil {
		return utils.BadRequest(c)
	}

	// get requesting user id
	targetId, ok := c.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(c)
	}

	notes, err := initializers.GrpcNoteService.GetAllNotes(c.UserContext(), &desc.GetAllNotesRequest{
		UserID:         targetId,
		Limit:          request.Limit,
		Offset:         request.Offset,
		AscendingOrder: request.AscendingOrder,
	})
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"notes": notes})
}

// GetNote godoc
// @Summary Get a specific note
// @Description Retrieves a specific note by ID for the authenticated user
// @Tags notes
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} map[string]interface{} "note"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /api/note/{id} [get]
func GetNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c)
	}

	// get user id
	targetId, ok := c.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(c)
	}

	// make request
	request := desc.GetNoteRequest{
		Id:     int64(id),
		UserID: targetId,
	}

	note, err := initializers.GrpcNoteService.GetNote(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"note": note})
}

// UpdateNote godoc
// @Summary Update a note
// @Description Updates an existing note for the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Param note body desc.UpdateNoteRequest true "Note update request"
// @Success 201 {object} map[string]interface{} "note"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /api/note [put]
func UpdateNote(c *fiber.Ctx) error {
	// get id
	targetId, ok := c.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(c)
	}

	// make request
	request := desc.UpdateNoteRequest{}
	if err := utils.ParseBodyAndValidate[desc.UpdateNoteRequest](c, &request); err != nil {
		return utils.BadRequest(c)
	}
	request.UserID = targetId

	// create note
	note, err := initializers.GrpcNoteService.UpdateNote(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"note": note})
}

// DeleteNote godoc
// @Summary Delete a note
// @Description Deletes a specific note by ID for the authenticated user
// @Tags notes
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} map[string]interface{} "deleted note"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /api/note/{id} [delete]
func DeleteNote(c *fiber.Ctx) error {
	// get user id
	targetId, ok := c.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(c)
	}

	// get note id
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c)
	}

	// make request
	request := desc.DeleteNoteRequest{
		Id:     int64(id),
		UserID: targetId,
	}

	_, err = initializers.GrpcNoteService.DeleteNote(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted note": id})
}
