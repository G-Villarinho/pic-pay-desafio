package utils

import (
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/go-playground/validator/v10"
)

func SetupCustomValidations(validator *validator.Validate) {
	validator.RegisterValidation("strongpassword", strongPasswordValidator)
	validator.RegisterValidation("cpf", cpfValidator)
}

func strongPasswordValidator(fl validator.FieldLevel) bool {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$&*])[A-Za-z\d!@#$&*]{8,}$`

	re := regexp2.MustCompile(pattern, 0)

	match, _ := re.MatchString(fl.Field().String())
	return match
}

func cpfValidator(fl validator.FieldLevel) bool {
	cpf := fl.Field().String()

	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	if len(cpf) != 11 {
		return false
	}

	if cpf == strings.Repeat(string(cpf[0]), 11) {
		return false
	}

	return isValidCPF(cpf)
}

func isValidCPF(cpf string) bool {
	var sum1, sum2, digit1, digit2 int
	for i := 0; i < 9; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum1 += num * (10 - i)
		sum2 += num * (11 - i)
	}

	digit1 = (sum1 * 10) % 11
	if digit1 == 10 {
		digit1 = 0
	}

	if digit1 != int(cpf[9]-'0') {
		return false
	}

	sum2 = digit1 * 2
	digit2 = (sum2 * 10) % 11
	if digit2 == 10 {
		digit2 = 0
	}

	return digit2 == int(cpf[10]-'0')
}
