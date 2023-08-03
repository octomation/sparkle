package diary

import "errors"

// Make sure you enabled "Daily notes" plugin and configured it.
// Suggest to check Config.Enabler() and Config.Section().
var configError = errors.New(
	"there is a problem with the Daily notes plugin configuration",
)

// Please review your "Daily notes" plugin configuration.
// Suggest to check Config.Section() with Config.FolderOptionPath().
var folderError = errors.New(
	"there is a problem with the Daily notes plugin folder",
)

// Please review your "Daily notes" plugin configuration.
// Suggest to check Config.Section() with Config.TemplateOptionPath().
var templateError = errors.New(
	"there is a problem with the Daily notes plugin template",
)
