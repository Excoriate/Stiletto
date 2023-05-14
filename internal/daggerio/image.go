package daggerio

var StackImagesMap = map[string]string{
	"PYTHON":           "python:3.8.5-slim-buster",
	"DOCKER":           "docker:23.0.1-dind",
	"INFRA:TERRAFORM":  "hashicorp/terraform",
	"INFRA:TERRAGRUNT": "alpine/terragrunt",
	//"AWS":        "amazon/aws-cli:latest",
	"AWS":    "alpine:latest",
	"ALPINE": "alpine:latest",
}

type DaggerContainerImage struct {
	Image   string
	Version string
}
