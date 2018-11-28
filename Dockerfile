FROM golang as builder

ARG GO_PACKAGE
RUN echo GO_PACKAGE=${GO_PACKAGE}

WORKDIR /go/src/${GO_PACKAGE}
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GO_OUTPUT_FILE=/app make

FROM scratch

COPY --from=builder /app /
ENV PATH /
STOPSIGNAL SIGTERM
ENTRYPOINT [ "app" ]