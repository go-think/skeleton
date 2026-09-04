package bootstrap

import (
	"app/app/providers"

	"github.com/go-think/think/support"
)

// AppProviders returns the application service providers.
func AppProviders() []support.ServiceProvider {
	return []support.ServiceProvider{
		&providers.AppServiceProvider{},
		&providers.EventServiceProvider{},
	}
}
