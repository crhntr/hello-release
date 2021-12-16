FROM pivotalcfreleng/kiln:0.63.0 as kiln
FROM bosh/cli2:6.4.8 as bosh

FROM golang:1.17-alpine
RUN go install github.com/pivotal-cf/om@7.4.1
COPY --from=bosh /usr/local/bin/bosh /usr/local/bin/bosh
COPY --from=kiln /kiln /usr/local/bin/kiln

WORKDIR /code