##################################
# STEP 1 build executable binary #
##################################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
COPY . $GOPATH/src/github.com/scotow/lingo

# Move to command directory.
WORKDIR $GOPATH/src/github.com/scotow/lingo/cmd/web

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/lingo

##############################
# STEP 2 build a small image #
##############################
FROM scratch

# Copy our static executable and static files.
COPY --from=builder /go/bin/lingo /lingo

# Run the hello binary.
ENTRYPOINT ["/lingo"]