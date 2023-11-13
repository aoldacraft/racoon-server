# Build the application from source
FROM docker.io/library/golang:alpine AS build-stage

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

#COPY go.mod go.sum ./
COPY . .
RUN go mod download

#COPY *.go ./

RUN go build -o /rc_server

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /rc_server /rc_server

EXPOSE 1323

ENTRYPOINT ["/rc_server"]

