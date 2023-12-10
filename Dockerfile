
FROM golang:1.21 AS build-stage
LABEL authors="Juan Daniel Sanchez Chavez"


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY jobs.toml /app/jobs.toml

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend

## Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /
ENV GIN_MODE=release

COPY --from=build-stage /app/backend /app
COPY --from=build-stage /app/jobs.toml /jobs.toml

EXPOSE 8081
EXPOSE 80
EXPOSE 443

USER nonroot:nonroot

ENTRYPOINT ["/app"]