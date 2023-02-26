# Execution Environment
FROM golang:1.20-alpine 

# Install Git.
RUN apk update && apk add --no-cache git

# Working Directory
RUN mkdir /build
WORKDIR /build

# Copying Files
COPY ./LogfileParser .

# Fetching dependencies.
RUN go mod download

# Set necessary environment variables needed
# for our image and build the sender.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o logfileparser.

# Building the binary executable.
RUN go build -o /logfileparser

# Running Executable
CMD [ "/logfileparser" ,"./user1.txt"]