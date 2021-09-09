FROM goreleaser/goreleaser:latest as builder

WORKDIR /go/src/github.com/benclapp/tesla_wall_connector_exporter

COPY go.* .
RUN go mod download
COPY . .

ARG BUILTBY
ARG COMMIT
ARG DATE
ARG VERSION

RUN GOARCH=amd64 GOOS=linux \
  go build \
  -ldflags "-X 'main.builtBy=${BUILTBY}' -X 'main.commit=${COMMIT}' -X 'main.date=${DATE}' -X 'main.version=${VERSION}'" \
  ./...

FROM scratch
COPY --from=builder \
  /go/src/github.com/benclapp/tesla_wall_connector_exporter/tesla_wall_connector_exporter \
  /tesla_wall_connector_exporter
ENTRYPOINT [ "/tesla_wall_connector_exporter" ]
