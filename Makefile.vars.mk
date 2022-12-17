## These are some common variables for Make

PROJECT_ROOT_DIR = .
PROJECT_NAME ?= paperless-cli
PROJECT_OWNER ?= ccremer

WORK_DIR = $(PWD)/.work

## BUILD:go
BIN_FILENAME ?= $(PROJECT_NAME)
go_bin ?= $(WORK_DIR)/bin
$(go_bin):
	@mkdir -p $@

## BUILD:docker
DOCKER_CMD ?= docker

IMG_TAG ?= latest
CONTAINER_REGISTRY ?= ghcr.io
# Image URL to use all building/pushing image targets
CONTAINER_IMG ?= $(CONTAINER_REGISTRY)/$(PROJECT_OWNER)/$(PROJECT_NAME):$(IMG_TAG)
