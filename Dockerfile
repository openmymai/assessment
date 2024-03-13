FROM golang:1.22.1-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -tags=unit ./...

RUN go build -o ./out/expenses .


### ------------

FROM alpine:3.19
COPY --from=build-base /app/out/expenses /app/expenses

ENV DATABASE_URL=postgres://ujngxgnh:Ejil-_AKm03rsFryOD4-vC_6BRM6jHS2@rain.db.elephantsql.com/ujngxgnh

ENV PORT=:2565

ENV AUTH_TOKEN="Basic YXBpZGVzaWduOjQ1Njc4"

EXPOSE 2565

CMD ["/app/expenses"]


# docker build -t expenses:multistage .
# docker run -i -t -p 2565:2565 expenses:multistage
# ❯ docker images expenses:multistage
# REPOSITORY   TAG          IMAGE ID       CREATED          SIZE
# expenses     multistage   f6f9c752bce5   37 seconds ago   16.7MB