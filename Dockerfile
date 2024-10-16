# Stage 2 - Your application image
FROM golang:1.22-alpine
WORKDIR /app

RUN apk update 
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache ffmpeg

COPY go.mod .
COPY go.sum .

COPY . .


# Download all dependencies
RUN go mod download
# Copy the source code into the container


 

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd
# Specify the command to run when the container starts

EXPOSE 8080


CMD [ "./main" ]