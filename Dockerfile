FROM alpine:3.22.2

RUN apk add --no-cache go

COPY . .

RUN go build .

EXPOSE 8080

CMD [ "./goParsingDocx" ]