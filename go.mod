module v2

go 1.13

replace github.com/kscloud/terraform-provider-dtekt/dtekt => ./dtekt

require github.com/hashicorp/terraform-plugin-sdk/v2 v2.3.0

require (
	github.com/hashicorp/terraform-plugin-test/v2 v2.0.0 // indirect
	github.com/kscloud/terraform-provider-dtekt/dtekt v1.0.0
)
