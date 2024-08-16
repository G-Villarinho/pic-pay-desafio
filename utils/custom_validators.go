package utils

import (
	"github.com/dlclark/regexp2"
	"github.com/klassmann/cpfcnpj"

	"github.com/go-playground/validator/v10"
)

func SetupCustomValidations(validator *validator.Validate) {
	validator.RegisterValidation("strongpassword", strongPasswordValidator)
	validator.RegisterValidation("cpf", cpfValidator)
	validator.RegisterValidation("uuid", uuidValidator)
}

func strongPasswordValidator(fl validator.FieldLevel) bool {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$&*])[A-Za-z\d!@#$&*]{8,}$`

	re := regexp2.MustCompile(pattern, 0)

	match, _ := re.MatchString(fl.Field().String())
	return match
}

func cpfValidator(fl validator.FieldLevel) bool {
	cpf := cpfcnpj.NewCPF(fl.Field().String())

	return cpf.IsValid()
}

func uuidValidator(fl validator.FieldLevel) bool {
	pattern := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`

	re := regexp2.MustCompile(pattern, 0)

	match, _ := re.MatchString(fl.Field().String())
	return match
}
