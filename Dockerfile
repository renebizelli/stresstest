FROM golang:latest AS builder

ENV url=x

WORKDIR /app

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o stresstest ./main.go

FROM scratch

COPY --from=builder /app/stresstest .

ENTRYPOINT [ "./stresstest", "stressOut"]
CMD []


#--url=www --requests=100 --concurrency=10
