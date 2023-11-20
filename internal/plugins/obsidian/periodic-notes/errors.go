package periodic

import "go.octolab.org/errors"

// Make sure you enabled "Periodic Notes" plugin and configured it.
// Suggest to check Config.Enabler() and Config.Section().
var errConfig errors.Message = "there is a problem with the Periodic Notes plugin configuration"

// Please review your "Periodic Notes" plugin configuration.
// Suggest to check Config.Section() for specified period.
var errFolder errors.Message = "there is a problem with the Daily notes plugin folder"

// Please review your "Periodic Notes" plugin configuration.
// Suggest to check Config.Section() for specified period.
var errTemplate errors.Message = "there is a problem with the Periodic Notes plugin template"
