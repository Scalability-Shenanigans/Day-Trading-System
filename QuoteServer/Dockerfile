FROM golang:1.20-alpine

# Install Git.
RUN apk update && apk add --no-cache git

# Working Directory
RUN mkdir /src
WORKDIR /src

# Copying Files
COPY . .

# Fetching dependencies.
RUN go get -d -v

# Building the binary executable.
RUN go build -o /quoteserver

# Running Executable
CMD [ "/quoteserver" ]