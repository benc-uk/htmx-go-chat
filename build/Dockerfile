# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./
COPY app/ ./app

RUN go mod download
RUN go build -o ./bin/server htmx-go-chat/app

# Final stage
FROM scratch
ARG VERSION=0.0.0

COPY --from=build /build/bin/server /server
COPY templates/ ./templates/

ENV VERSION=$VERSION
ENV PORT=8000

EXPOSE 8000

CMD ["/server"]
