# =====================================================
# Build stage
# =====================================================
FROM golang:1.21-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY app/ ./app
RUN go build -o ./server htmx-go-chat/app

# =====================================================
# Main runtime stage
# =====================================================
FROM scratch
ARG VERSION=0.0.0

COPY --from=build /build/server /server
COPY templates/ ./templates/

VOLUME /data

ENV VERSION=$VERSION
ENV PORT=8000
ENV DB_PATH=/data

EXPOSE 8000

CMD ["/server"]
