package validatorha

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

// IValidator defines the methods for validation
type IValidator interface {
	AddRuleStruct(interface{}) error
	ValidateStruct(interface{}) error
}

// ValidatorConfig holds configuration for the validator
type ValidatorConfig struct {
	EnableIPv4Validation bool                      // Flag to enable/disable IPv4 validation
	CustomMessages       map[string]string         // Custom error messages for validations
	CustomValidators     map[string]validator.Func // Custom validation functions
}

// Validator implements the IValidator interface
type Validator struct {
	config   ValidatorConfig
	validate *validator.Validate
}

// NewValidator creates a new Validator with the given config
func NewValidator(config ValidatorConfig) *Validator {
	v := validator.New()

	// Register custom validations if enabled in config

	if config.EnableIPv4Validation {
		v.RegisterValidation("ipv4", validateIPv4)
	}

	// Register any custom validators provided in the config
	for tag, fn := range config.CustomValidators {
		v.RegisterValidation(tag, fn)
	}

	// Register custom error messages if provided
	if config.CustomMessages != nil {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "-" {
				return ""
			}
			return name
		})
	}

	return &Validator{
		config:   config,
		validate: v,
	}
}

// AddRuleStruct is used to add custom validation rules
func (v *Validator) AddRuleStruct(rule interface{}) error {
	// Simple function implementation, assuming no additional logic needed for now
	return nil
}

// ValidateStruct validates a struct based on the registered rules

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// validateIPv4 is a custom validation function for IPv4 addresses
func validateIPv4(fl validator.FieldLevel) bool {
	ip := fl.Field().String()
	return validator.New().Var(ip, "ipv4") == nil
}
