package config

// LimitByNumOfFilesChanged compat function until codebase is migrated to use the new config system
func LimitByNumOfFilesChanged() bool {
	// if this flag is set then it will fail if there are more projects impacted than the
	// number of files changed
	return BackendConfig.Features.LimitMaxProjectsToFilesChanged
}

// GetPort compat function until codebase is migrated to use the new config system
func GetPort() int {
	return BackendConfig.Server.Port
}
