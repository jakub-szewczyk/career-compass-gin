# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.1
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./
RUN go mod download -x

# Copy the entire source code including templates
COPY . .

# Build the Go binary
ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

################################################################################
# Create the final minimal runtime image
FROM alpine:latest AS final

# Install runtime dependencies
RUN apk --update add ca-certificates tzdata && update-ca-certificates

# Create a non-privileged user
ARG UID=10001
RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "${UID}" appuser

# Switch to non-root user
USER appuser

# Copy the built server binary
COPY --from=build /bin/server /bin/server

# Copy the templates directory explicitly
COPY --from=build /src/templates /templates

# Set correct permissions
USER root
RUN chown -R appuser:appuser /templates && chmod -R 755 /templates
USER appuser

# Expose the port
EXPOSE 3000

# Run the server
ENTRYPOINT ["/bin/server"]
