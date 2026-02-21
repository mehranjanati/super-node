package wazero

import "go.uber.org/fx"

// Module is the Fx module for the Wazero service.
var Module = fx.Module("wazero",
	fx.Provide(NewWazeroService),
)
