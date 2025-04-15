FROM golang:1.24.1 as builder

ARG COMMIT_SHA
ARG VERSION="v0.0.0"
ARG BUILD_DATE

RUN echo "Building version: ${VERSION}, commit: ${COMMIT_SHA}, date: ${BUILD_DATE}"

# Set the working directory
WORKDIR /src

# Copy all required source, blacklist files that are not required through `.dockerignore`
COPY . .

# Build with optimizations and embedded version info
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags="-X 'github.com/aldoborrero/strate/backend/version.Version=${VERSION}' \
              -X 'github.com/aldoborrero/strate/backend/version.Meta=${COMMIT_SHA}' \
              -X 'github.com/aldoborrero/strate/backend/version.BuildDate=${BUILD_DATE}' \
              -s -w" \
    -o st-server ./backend/

# Ensure the binary is executable
RUN chmod +x /src/st-server

# Runtime stage using distroless
FROM gcr.io/distroless/base-debian12:nonroot as runner

LABEL org.opencontainers.image.source="https://github.com/aldoborrero/strate"
LABEL org.opencontainers.image.description="Strate Backend Server"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT_SHA}"
LABEL org.opencontainers.image.created="${BUILD_DATE}"

# Set production mode
ENV GIN_MODE=release

# Expose the running port
EXPOSE 3000

# Copy required files
WORKDIR /app
COPY --from=builder /src/st-server /app/st-server
COPY --from=builder /src/backend/migrations /app/migrations

# Use the nonroot user that comes with the distroless image
USER nonroot:nonroot

# Define the entrypoint
ENTRYPOINT ["/app/st-server"]
