FROM golang as builder

ARG package
WORKDIR /go/src/${package}
ADD . .
RUN CGO_ENABLED=0 GOOS=linux OUTPUT=/app make

FROM scratch

COPY --from=builder /app /
ENV PATH /
STOPSIGNAL SIGTERM
ENTRYPOINT [ "app" ]