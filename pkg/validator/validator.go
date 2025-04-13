package validator

import (
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/yogarn/filkompedia-be/model"
)

func RegisterValidator(v *validator.Validate) {
	v.RegisterValidation("rfc3339date", Date)
	v.RegisterValidation("image_type", ImageType)
	v.RegisterValidation("image_size", ImageSize)
}

func Date(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

func ImageType(fl validator.FieldLevel) bool {
	fileHeader, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}

	mimeType, err := model.GetImageType(&fileHeader)
	if err != nil {
		return false
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	return allowedTypes[mimeType]
}

func ImageSize(fl validator.FieldLevel) bool {
	fileHeader, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}

	const maxSize = 2 * 1024 * 1024
	return fileHeader.Size <= maxSize
}
