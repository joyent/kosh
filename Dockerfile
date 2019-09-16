# vim: se syn=dockerfile:
FROM golang:1.13-alpine AS build
ENV CGO_ENABLED 0

RUN apk add --no-cache --update make git perl-utils dep shadow

ENV PATH "/go/bin:${PATH}"

RUN go get honnef.co/go/tools/cmd/staticcheck

RUN mkdir -p /go/src/github.com/joyent/kosh
WORKDIR /go/src/github.com/joyent/kosh

COPY . /go/src/github.com/joyent/kosh/

RUN make

FROM scratch
COPY --from=build /go/src/github.com/joyent/kosh/bin/kosh /bin/kosh
COPY --from=build /etc/ssl /etc/ssl

ENV KOSH_TOKEN "broken"
ENTRYPOINT [ "/bin/kosh" ]
CMD ["version"]
