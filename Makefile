PROJECT_NAME ?= webzr
DESCRIPTION ?= crappy websocket research

# Go gonfig
GO ?= go
GOMCD = $(shell which go)
GOFILES = $(shell find . -type f -name '*.go' -not -path "./.git/*")


