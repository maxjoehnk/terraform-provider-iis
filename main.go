package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/maxjoehnk/terraform-provider-iis/iis"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: iis.Provider,
	})
}
