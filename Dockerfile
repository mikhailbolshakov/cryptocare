FROM IMAGE

ENV GO111MODULE=on

WORKDIR /src

COPY . ./

#RUN make swagger

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN ls -l && pwd

RUN make build

FROM alpine:3.12.3

RUN apk --no-cache add ca-certificates

ARG SERVICE

WORKDIR /usr/local/root/trading

ENV CRYPTOCAREROOT="/usr/local/root"

RUN ls -la

COPY --from=0 /src/bin/main ./bin/main
COPY --from=0 /src/config.yml ./config.yml
COPY ./src/db/migrations ./src/db/migrations

RUN ls -la && pwd

ENTRYPOINT ["./bin/main"]
