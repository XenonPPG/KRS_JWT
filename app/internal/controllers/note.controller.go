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
	return utils.GrpcHandler(c, initializers.GrpcNoteService.CreateNote)
}

func GetAllNotes(c *fiber.Ctx) error {
	request := models.GetAllItemsRequest{}

	// parse query
	if err := c.QueryParser(&request); err != nil {
		return utils.BadRequest(c)
	}

	notes, err := initializers.GrpcNoteService.GetAllNotes(c.UserContext(), &desc.GetAllNotesRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"notes": notes})
}

func GetNote(c *fiber.Ctx) error {
	request := desc.GetNoteRequest{}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c)
	}
	request.Id = int64(id)

	note, err := initializers.GrpcNoteService.GetNote(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"note": note})
}

func UpdateNote(c *fiber.Ctx) error {
	return utils.GrpcHandler(c, initializers.GrpcNoteService.UpdateNote)
}

func DeleteNote(c *fiber.Ctx) error {
	request := desc.DeleteNoteRequest{}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c)
	}
	request.Id = int64(id)

	_, err = initializers.GrpcNoteService.DeleteNote(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted note": id})
}
