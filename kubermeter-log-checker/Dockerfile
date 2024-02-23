FROM golang:1.19-alpine AS build
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}
ENV HOME /app
ADD . /build/
WORKDIR /build
RUN chown -R 777 /build && go build -o kube-stopper main.go 

FROM alpine AS runtime

ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}
ENV HOME /app
RUN mkdir -p /app \
  && adduser -S -D -H -h /app appuser
COPY --from=build /build/kube-stopper ./app/
RUN chown -R appuser /app
WORKDIR /app
USER appuser
EXPOSE 8080
CMD ["./kube-stopper"]

