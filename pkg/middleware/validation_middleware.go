package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// ValidateStruct валидирует структуру используя теги valid
func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}

// ValidationMiddleware создает middleware для автоматической валидации тела запроса
func ValidationMiddleware(modelType interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(modelType); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "Failed to parse request body",
				"details": err.Error(),
			})
		}

		if err := ValidateStruct(modelType); err != nil {
			validationErrors := make([]string, 0)

			if _, ok := err.(*validator.InvalidValidationError); ok {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": "Validation error",
				})
			}

			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, err.Error())
			}

			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Validation Failed",
				"message": "Request validation failed",
				"details": validationErrors,
			})
		}

		return c.Next()
	}
}
