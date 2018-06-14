FROM alpine:3.7

EXPOSE 8080

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH 

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} \
     && adduser -u 1001 -S -G root -g 0 -D -h ${APP_ROOT} -s /sbin/nologin go 

COPY ./hicicd ${APP_BIN}

# Drop the root user and make the content of /opt/app-root owned by user 1001
RUN chown -R 1001:0 ${APP_ROOT}

RUN chmod 755 ${APP_BIN}/hicicd

WORKDIR ${APP_ROOT}

USER 1001

CMD [ "hicicd" ]
