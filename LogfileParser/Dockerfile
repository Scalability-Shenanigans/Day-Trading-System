# Execution Environment
FROM golang:1.16-alpine

# Working Directory
WORKDIR /

# Copying Files
COPY *.go ./
COPY go.mod ./

# Build 
RUN go build -o /logfileparser

# Expose Port 8080
EXPOSE 8080

# Running Executable
CMD [ "./logfileparser" ]