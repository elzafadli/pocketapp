# Stage 1 from golang v1.21 base image
FROM golang:1.25.5 as builder

# Add Maintainer Info
LABEL maintainer="devops"

# Set the Current Working Directory inside the container
WORKDIR /go/src/app
ADD .  /go/src/app

# Ensure version.txt exists
RUN if [ ! -f version.txt ]; then touch version.txt; fi

# Build project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./pakuningratan



# # Stage 2
FROM alpine:3

# Add Maintainer Info
LABEL maintainer="devops"

RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy binary from builder
COPY --from=builder /go/src/app/config.yaml /app/config.yaml
COPY --from=builder /go/src/app/.env /app/.env
COPY --from=builder /go/src/app/pakuningratan /app/pakuningratan
COPY --from=builder /go/src/app/version.txt /app/version.txt

# Run the executable
CMD ["./main","service"]
