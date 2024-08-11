FROM alpine:3.18
WORKDIR /scripts

RUN apk update && apk add bash postgresql15-client wget

RUN wget -O /usr/bin/goose https://github.com/pressly/goose/releases/download/v3.7.0/goose_linux_x86_64 && \
	chmod a+x /usr/bin/goose && \
	chown postgres:postgres /usr/bin/goose

# local
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=erp_db
ENV DB_USER_RW_PASSWORD=postgres

# remote
# ENV DB_HOST=31.131.254.247
# ENV DB_PORT=5432
# ENV DB_USER=erp_db_usr
# ENV DB_PASSWORD=E4SeIFlKFSCRBQIz
# ENV DB_NAME=erp_db

ADD migrations/data_base/*.sql data_base/
ADD migrations/database-init-up.sh .
COPY . .
RUN chmod +x database-init-up.sh

ENTRYPOINT [ "./database-init-up.sh" ]


