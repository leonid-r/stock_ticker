FROM umputun/baseimage:buildgo as build

WORKDIR /build
ADD . /build

RUN golangci-lint run --out-format=tab --tests=false ./...

RUN \
    revison=$(/script/git-rev.sh) && \
    echo "revision=${revison}" && \
    go build -mod=vendor -o app -ldflags "-X main.revision=$revison -s -w" .


FROM umputun/baseimage:app

COPY --from=build /build/app /srv/app

EXPOSE 8080
WORKDIR /srv

CMD ["/srv/app"]