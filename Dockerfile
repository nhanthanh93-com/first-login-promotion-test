FROM golang:1.21-alpine

# Set the current working directory inside the container
WORKDIR /go/src/ta_app
ARG env
ARG config
ARG version
ENV env=${env}
ENV version=${version}
ENV config=${config}
ENV GIN_MODE=release

#RUN setcap CAP_NET_BIND_SERVICE=+eip /app/app-exe
# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod tidy
RUN go install github.com/cosmtrek/air@v1.45.0
RUN go mod download

# Copy the source from the current directory to the workspace
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Command to run the executable
CMD ["air"]