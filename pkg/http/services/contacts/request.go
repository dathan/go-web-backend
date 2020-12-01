package contacts

import "github.com/System-Glitch/goyave/v3/validation"

var (
	AddContactRequest validation.RuleSet = validation.RuleSet{
		"email": {"required", "string"},
	}
)
