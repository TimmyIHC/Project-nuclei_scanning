# Use the latest Go parent image
FROM golang:latest

# Set the working directory inside the container
RUN mkdir /app
WORKDIR /app

# Installing nuclei
RUN go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest

RUN go mod init requirements

RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/s3/s3manager

# Copy your Go source code into the container
COPY go_entrypoint.go .

# Define the command to run the Go application
ENTRYPOINT ["go", "run"]
CMD ["go_entrypoint.go"]