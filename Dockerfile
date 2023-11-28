# Start from golang alpine base image
FROM golang:1.19-alpine
LABEL authors="Juan Daniel Sanchez Chavez"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code to the Working Directory inside the container
COPY . .

# Download all dependencies.
# Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o Backend main.go

# Export necessary ports
EXPOSE 8080

# Command to run when starting the container
CMD ["/app/job-app"]
