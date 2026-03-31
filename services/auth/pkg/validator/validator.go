package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func ValidatePayload(payload any) error {

	if err := Validate.Struct(payload); err != nil {

		var verr validator.ValidationErrors

		var fieldsTags []string

		if errors.As(err, &verr) {
			for _, v := range verr {
				msg := fmt.Sprintf(
					"'%s' failed on the '%s' tag",
					v.Field(),
					v.Tag(),
				)
				fieldsTags = append(fieldsTags, msg)
			}

			// v := verr[0]
			return errors.New(strings.Join(fieldsTags, ","))
		}

		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
