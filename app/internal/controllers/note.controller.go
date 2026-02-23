package controllers

import (
	"JWT/internal/initializers"
	"JWT/internal/models"
	"JWT/internal/utils"
	"strconv"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/note_v1"
	"github.com/gofiber/fiber/v2"
)

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
		UserID: targetId,
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"notes": notes})
}

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
