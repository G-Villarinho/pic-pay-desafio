package domain

import (
	"github.com/dlclark/regexp2"
	"github.com/google/uuid"
	"github.com/klassmann/cpfcnpj"

	"github.com/go-playground/validator/v10"
)

const (
	StrongPasswordTag = "strongpassword"
	CPFTag            = "cpf"
	UUIDTag           = "uuid"
	WalletTypeTag     = "wallettype"
)

func SetupCustomValidations(validator *validator.Validate) {
	validator.RegisterValidation("strongpassword", strongPasswordValidator)
	validator.RegisterValidation("cpf", cpfValidator)
	validator.RegisterValidation("uuid", uuidValidator)
	validator.RegisterValidation("wallettype", walletTypeValidator)
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
	_, err := uuid.Parse(fl.Field().String())

	return err != nil
}

func walletTypeValidator(fl validator.FieldLevel) bool {
	walletType, ok := fl.Field().Interface().(WalletType)
	if !ok {
		return false
	}
	return walletType.IsValid()
}
