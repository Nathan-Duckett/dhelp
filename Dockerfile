FROM golang:1.21 AS dhelp-build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /dhelp

# Run the tests in the container
FROM dhelp-build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=dhelp-build-stage /dhelp /dhelp

EXPOSE 4000

ENTRYPOINT ["/dhelp"]