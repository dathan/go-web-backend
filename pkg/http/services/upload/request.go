package upload

import "github.com/System-Glitch/goyave/v3/validation"

var (
	UploadRequest validation.RuleSet = validation.RuleSet{
		"file": {"file", "mime:text/csv"},
	}
)
