# pull the base image

FROM golang:latest AS builder
# create base working directory inside container
WORKDIR /app
# It is a best practice to copy these files separately
# before copying the rest of the project files because
# the dependencies are unlikely to change as frequently
# as the source code. This reduces the number of layers
# Docker needs to rebuild if only the source code has
# changed, and speeds up the build process.
COPY go.mod go.sum ./
# Install all the dependencies
RUN go mod tidy
# Now copy all
COPY . .
# build the go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-bookserver .

# Now this image takes way too much space because of the base image
# We will use multi staging to reduce the space used.

# pull alpine latest image, it contains necessary files for running go binaries
FROM alpine:latest
WORKDIR /app
# COPY everything from last image working directory
COPY --from=builder /app/api-bookserver .
ENTRYPOINT ["./api-bookserver"]
# Set default command when a container will be started
# In my case, my start command starts the server in 8081 port by default
CMD ["start"]