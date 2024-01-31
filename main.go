package main

import (
	"github.com/Pradipbabar/todo/provider"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *provider.Provider {
			return provider.Provider()
		},
	})
}
