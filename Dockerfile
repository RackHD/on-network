# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /project
WORKDIR /project

# Build the linux binary for on-network
RUN make link clean deps linux

# Run the on-network command by default when the container starts.
CMD /project/cmd/on-network-server/on-network-linux-amd64 --port 8080 --host 0.0.0.0 --write-timeout 10m

# Document that the service listens on port 8080.
EXPOSE 8080
