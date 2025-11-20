package services

import (
    "regexp"
    "strings"
    "unicode"

    "github.com/arman300s/uni-portal/internal/core/contracts"
)

const (
    minPasswordLength = 8
    maxPasswordLength = 128
    maxNameLength     = 100
    maxEmailLength    = 255
)

func validateSignupInput(input contracts.SignupInput) contracts.ValidationErrors {
    var errs contracts.ValidationErrors
    if err := validateName(input.Name); err != nil {
        errs = append(errs, extractValidationErrors(err, "name")...)
    }
    if err := validateEmail(input.Email); err != nil {
        errs = append(errs, extractValidationErrors(err, "email")...)
    }
    if err := validatePassword(input.Password); err != nil {
        errs = append(errs, extractValidationErrors(err, "password")...)
    }
    return errs
}

func validateCreateUserInput(input contracts.CreateUserInput) contracts.ValidationErrors {
    var errs contracts.ValidationErrors
    if err := validateName(input.Name); err != nil {
        errs = append(errs, extractValidationErrors(err, "name")...)
    }
    if err := validateEmail(input.Email); err != nil {
        errs = append(errs, extractValidationErrors(err, "email")...)
    }
    if err := validatePassword(input.Password); err != nil {
        errs = append(errs, extractValidationErrors(err, "password")...)
    }
    if strings.TrimSpace(input.RoleName) == "" {
        errs = append(errs, contracts.ValidationError{Field: "role", Message: "role is required"})
    }
    return errs
}

func validateSubjectInput(input contracts.SubjectInput) contracts.ValidationErrors {
    var errs contracts.ValidationErrors
    if strings.TrimSpace(input.Name) == "" {
        errs = append(errs, contracts.ValidationError{Field: "name", Message: "name is required"})
    }
    return errs
}

func validateLoginInput(input contracts.LoginInput) contracts.ValidationErrors {
    var errs contracts.ValidationErrors
    if strings.TrimSpace(input.Email) == "" {
        errs = append(errs, contracts.ValidationError{Field: "email", Message: "email is required"})
    }
    if strings.TrimSpace(input.Password) == "" {
        errs = append(errs, contracts.ValidationError{Field: "password", Message: "password is required"})
    }
    return errs
}


func validateUpdateUserInput(input contracts.UpdateUserInput) contracts.ValidationErrors {
    var errs contracts.ValidationErrors
    if trimmed := strings.TrimSpace(input.Email); trimmed != "" {
        if err := validateEmail(trimmed); err != nil {
            errs = append(errs, extractValidationErrors(err, "email")...)
        }
    }
    if strings.TrimSpace(input.RoleName) == "" {
        errs = append(errs, contracts.ValidationError{Field: "role", Message: "role is required"})
    }
    return errs
}

func validateEmail(email string) error {
    trimmed := strings.TrimSpace(strings.ToLower(email))
    switch {
    case trimmed == "":
        return newValidationMsg("email", "email is required")
    case len(trimmed) > maxEmailLength:
        return newValidationMsg("email", "email is too long")
    }
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(trimmed) {
        return newValidationMsg("email", "invalid email format")
    }
    return nil
}

func validatePassword(password string) error {
    switch {
    case password == "":
        return newValidationMsg("password", "password is required")
    case len(password) < minPasswordLength:
        return newValidationMsg("password", "password must be at least 8 characters")
    case len(password) > maxPasswordLength:
        return newValidationMsg("password", "password is too long")
    }
    var (
        hasUpper   bool
        hasLower   bool
        hasNumber  bool
        hasSpecial bool
    )
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    switch {
    case !hasUpper:
        return newValidationMsg("password", "password must contain at least one uppercase letter")
    case !hasLower:
        return newValidationMsg("password", "password must contain at least one lowercase letter")
    case !hasNumber:
        return newValidationMsg("password", "password must contain at least one number")
    case !hasSpecial:
        return newValidationMsg("password", "password must contain at least one special character")
    }
    return nil
}

func validateName(name string) error {
    trimmed := strings.TrimSpace(name)
    switch {
    case trimmed == "":
        return newValidationMsg("name", "name is required")
    case len(trimmed) > maxNameLength:
        return newValidationMsg("name", "name is too long")
    }
    nameRegex := regexp.MustCompile(`^[a-zA-Z\\s\\-']+$`)
    if !nameRegex.MatchString(trimmed) {
        return newValidationMsg("name", "name contains invalid characters")
    }
    return nil
}

func newValidationMsg(field, message string) error {
    return contracts.ValidationErrors{contracts.ValidationError{Field: field, Message: message}}
}

func extractValidationErrors(err error, fallbackField string) contracts.ValidationErrors {
    switch v := err.(type) {
    case contracts.ValidationErrors:
        return v
    default:
        return contracts.ValidationErrors{contracts.ValidationError{
            Field:   fallbackField,
            Message: err.Error(),
        }}
    }
}
