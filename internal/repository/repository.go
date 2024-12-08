package repository

import (
	"career_backend/internal/dal"
	"github.com/google/wire"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(dal.Init, NewUserRepo, NewCompanyRepo)
