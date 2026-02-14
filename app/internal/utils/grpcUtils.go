package utils

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcHandler[Request any, Response any](
	c *fiber.Ctx,
	call func(context.Context, *Request, ...grpc.CallOption) (*Response, error)) error {
	req := new(Request)
	if err := ParseBodyAndValidate[Request](c, req); err != nil {
		return err
	}

	res, err := call(c.UserContext(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return InternalServerError(c)
		}

		switch st.Code() {
		case codes.InvalidArgument:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": st.Message()})
		case codes.NotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"err": st.Message()})
		case codes.AlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"err": st.Message()})
		case codes.Unauthenticated:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"err": st.Message()})
		case codes.PermissionDenied:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"err": st.Message()})
		case codes.Unavailable:
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"err": st.Message()})
		case codes.DeadlineExceeded:
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"err": st.Message()})
		default:
			return InternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
