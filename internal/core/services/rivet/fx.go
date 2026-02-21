package rivet

import "go.uber.org/fx"

// Module is the Fx module for the Rivet service.
var Module = fx.Module("rivet",
	fx.Provide(NewRivetService),
)
