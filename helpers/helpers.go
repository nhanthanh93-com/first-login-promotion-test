package helpers

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func Validate(str interface{}) error {
	validate := validator.New()
	err := validate.Struct(str)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUUID(uuid string) bool {
	regex := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(uuid)
}
