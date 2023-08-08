# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.245.2/containers/go/.devcontainer/base.Dockerfile

ARG VARIANT="1.20"
ARG LINT="1.53.3"

FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# Update packages
RUN apt-get update

# [Package] Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b $(go env GOPATH)/bin v1.53.3

# [Package] Install act
RUN curl -s https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
