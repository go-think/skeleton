package providers

import (
	"app/app/services"
	"app/config"

	"github.com/go-think/think/container"
	"github.com/go-think/think/contract"
)

// AppServiceProvider registers application-level services and configurations.
type AppServiceProvider struct{}

// Register registers container bindings and business configurations.
func (p *AppServiceProvider) Register(app *container.Container) {
	// 1. Merge business configurations into the framework configuration repository
	if cfg := app.Make[contract.Config]("contract.Config"); cfg != nil {
		cfg.Load(config.All())
	}

	// 2. Singleton binding
	app.Singleton[*services.UserService](func() *services.UserService {
		return services.NewUserService()
	}, "UserService")
	app.Singleton[*services.UserService](func() *services.UserService {
		return services.NewUserService()
	})

	// 3. Transient binding
	app.Bind[*services.OrderGenerator](func() *services.OrderGenerator {
		return services.NewOrderGenerator()
	}, "OrderGenerator")

	// 4. Scoped binding (per-request)
	app.Scoped[*services.RequestContext](func() *services.RequestContext {
		return services.NewRequestContext()
	}, "RequestContext")

	// 5. Tagged bindings
	app.Singleton[services.PaymentService](func() services.PaymentService {
		return &services.WechatPay{}
	}, "pay.wechat")
	app.Singleton[services.PaymentService](func() services.PaymentService {
		return &services.AliPay{}
	}, "pay.alipay")
	app.Tag("pay.wechat", "payments")
	app.Tag("pay.alipay", "payments")
}

// Boot bootstraps application services
func (p *AppServiceProvider) Boot(app *container.Container) {
}
