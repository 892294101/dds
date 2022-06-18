GOCMD=go
GOBUILD=${GOCMD} build -gcflags=all='-l -N' -ldflags '-s -w'

BUILD_DIR=../../build
BINARY_DIR=$(BUILD_DIR)/bin

MYSQL_EXTRACT_FILE=$(BINARY_DIR)/mysqlextract
MYSQL_EXTRACT_SRC=./mysqlextract.go


.PHONY: all clean build

all: clean build

clean:

build:
	${GOBUILD} -o ${MYSQL_EXTRACT_FILE} ${MYSQL_EXTRACT_SRC}

