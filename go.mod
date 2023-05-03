module github.com/cvbarros/terraform-provider-teamcity

go 1.13

require (
	github.com/cvbarros/go-teamcity v1.1.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.26.1
	github.com/motemen/go-nuts v0.0.0-20190725124253-1d2432db96b0 // indirect
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect
)

replace github.com/cvbarros/go-teamcity => ../go-teamcity
