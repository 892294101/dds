GOCMD=go
GOBUILD=${GOCMD} build -gcflags=all='-l -N' -ldflags '-s -w'

BUILD_DIR=../build
BINARY_DIR=$(BUILD_DIR)/bin
PLUGIN_DIR=$(BUILD_DIR)/lib
LOGS_DIR=$(BUILD_DIR)/logs

BINARY_FILE=$(BINARY_DIR)/cli
BINARY_SRC=./cli.go

PLUGINS_SYSCMD_FILE=$(PLUGIN_DIR)/sys_command.so
PLUGINS_SYSCMD_SRC=./terminal/plugins/syscmd.go

PLUGINS_EXTRACTCMD_FILE=$(PLUGIN_DIR)/extract_command.so
PLUGINS_EXTRACTCMD_SRC=./terminal/plugins/extractcmd.go

.PHONY: all clean build

all: clean build

clean:
	@if [ -f  ${BUILD_DIR} ]; then rm -rf ${BUILD_DIR}/*; fi
	mkdir -p ${BINARY_DIR} ${PLUGIN_DIR} ${LOGS_DIR}

build:
	${GOBUILD} -o ${BINARY_FILE} ${BINARY_SRC}
	${GOBUILD} -buildmode=plugin -o ${PLUGINS_SYSCMD_FILE} ${PLUGINS_SYSCMD_SRC}
	${GOBUILD} -buildmode=plugin -o ${PLUGINS_EXTRACTCMD_FILE} ${PLUGINS_EXTRACTCMD_SRC}

