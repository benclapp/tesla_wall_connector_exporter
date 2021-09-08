FROM goreleaser/goreleaser:latest as builder

WORKDIR /go/src/github.com/benclapp/tesla_wall_connector_exporter
COPY . .
# Remove snapshot flag on merge to main
RUN GOARCH=amd64 goreleaser build --single-target --snapshot

FROM scratch
COPY --from=builder \
  /go/src/github.com/benclapp/tesla_wall_connector_exporter/dist/tesla_wall_connector_exporter_linux_amd64/tesla_wall_connector_exporter \
  /tesla_wall_connector_exporter
ENTRYPOINT [ "/tesla_wall_connector_exporter" ]
