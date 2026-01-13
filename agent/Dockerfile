FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o agent ./cmd

FROM alpine
WORKDIR /app
COPY --from=build /app/agent .
EXPOSE 9090
CMD ["./agent"]
