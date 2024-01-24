# Build stage
FROM golang:1.21.0-alpine3.17 AS build
LABEL "dev.aliaktas"="AA"
LABEL mantainer="ali@aliaktas.dev"

WORKDIR /app
COPY . /app
RUN go build -o herald

# Run stage
FROM alpine:3.18.3
WORKDIR /app
COPY --from=build /app/herald .
COPY --from=build /app/template ./template

CMD ["./herald"]