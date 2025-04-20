FROM alpine:latest

RUN apk update --no-cache

RUN apk add --no-cache postgresql-client bash

COPY /database/postgres/init_schemas/init-schemas.sh /

RUN chmod +x /init-schemas.sh

ENTRYPOINT [ "./init-schemas.sh" ]