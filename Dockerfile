FROM golang:1.22.5-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o app .

RUN mkdir -p /mnt/img_orig /mnt/img_res

EXPOSE 8085

CMD ["./app", "-path-orig", "/mnt/img_orig", "-path-res", "/mnt/img_res"]
