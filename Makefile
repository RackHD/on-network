# Borrowed from:
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = on-network
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=RackHD
GITHUB_ORG=${GOPATH}/src/github.com/${GITHUB_USERNAME}
REPO_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
BUILD_DIR=${REPO_DIR}/cmd/${BINARY}-server
CURRENT_DIR=$(shell pwd)
REPO_DIR_LINK=$(shell readlink ${REPO_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: link clean deps test-deps test vet linux darwin windows

deps:
	go get github.com/Masterminds/glide
	glide install

test-deps:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/modocache/gover
	go get github.com/mattn/goveralls

link:
	REPO_DIR=${REPO_DIR}; \
	REPO_DIR_LINK=${REPO_DIR_LINK}; \
	CURRENT_DIR=${CURRENT_DIR}; \
	GITHUB_ORG=${GITHUB_ORG}; \
	if [ "$${REPO_DIR_LINK}" != "$${CURRENT_DIR}" ]; then \
	    echo "Fixing symlinks for build"; \
	    rm -f $${REPO_DIR}; \
	    mkdir -p $${GITHUB_ORG}; \
	    ln -s $${CURRENT_DIR} $${REPO_DIR}; \
	fi

linux:
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} . ; 

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; 

windows:
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; 

docker:
	docker build -t on-network:dev .
test:
	cd ${REPO_DIR}; \
	ginkgo -v ./... -cover -trace -race ; \
	gover ... ; \
	cat gover.coverprofile ; 

vet:
	cd ${REPO_DIR}; \
	go vet ./... > ${VET_REPORT} 2>&1 ; 

fmt:
	cd ${REPO_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; 

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

.PHONY: link linux darwin windows test vet fmt clean
