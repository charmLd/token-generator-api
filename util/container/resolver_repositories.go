package container

import (
	"github.com/charmLd/token-generator-api/externals/repositories"
	"time"
)

var resolvedRepositories Repositories

// Resolve all repositories.
func resolveRepositories(location *time.Location) Repositories {
	resolveUserRepository(location)

	resolveTokenRepository(location)

	return resolvedRepositories
}

// Resolve the user repository
func resolveUserRepository(location *time.Location) {

	userRepo := repositories.NewUserRepository(resolvedAdapters.MySQL, location)
	resolvedRepositories.UserRepository = userRepo
}

// Resolve token repository
func resolveTokenRepository(location *time.Location) {
	repo := repositories.NewTokenRepository(resolvedAdapters.MySQL, location)
	resolvedRepositories.TokenRepository = repo
}
