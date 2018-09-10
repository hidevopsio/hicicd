FROM alpine:3.8

EXPOSE 8080
EXPOSE 7575

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH 

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} \
     && apk update \
     && apk upgrade \
     && apk --no-cache add ca-certificates\
     && apk add -U tzdata ttf-dejavu\
     && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
     && adduser -u 1001 -S -G root -g 0 -D -h ${APP_ROOT} -s /sbin/nologin go

COPY ./hicicd ${APP_BIN}
COPY ./config ${APP_ROOT}/config
# Drop the root user and make the content of /opt/app-root owned by user 1001
RUN chown -R 1001:0 ${APP_ROOT}

RUN chmod 755 ${APP_BIN}/hicicd

WORKDIR ${APP_ROOT}

USER 1001

CMD [ "hicicd" ]
