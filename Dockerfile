# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20.3 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o travelbot ./cmd/botservice/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/travelbot /travelbot

EXPOSE 2112

USER nonroot:nonroot # user 65532:65532 is nonroot in distroless

ENTRYPOINT ["/travelbot"]
CMD ["--sleep=60"]