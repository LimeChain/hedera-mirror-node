# -----------------------------  Rosetta  ----------------------------- #
FROM golang:1.13 as rosetta-builder
WORKDIR /tmp/src/hedera-mirror-rosetta
COPY ./hedera-mirror-rosetta . 
RUN go build -o rosetta-executable ./cmd
# -----------------------------  Rosetta END ----------------------------- #

# ----------------------------- Importer/GRPC ---------------------------- #
FROM openjdk:11.0 as java-builder

RUN apt-get update && apt-get install -y git
RUN git clone https://github.com/LimeChain/hedera-mirror-node.git
# COPY ./application.yml ./hedera-mirror-node/application.yml
RUN cd hedera-mirror-node && ./mvnw clean package -DskipTests

# -------------------------- Importer/GRPC END --------------------------- #
FROM ubuntu:16.04 as runner

# ----------------------------- PosgreSQL ----------------------------- #
# Add the PostgreSQL PGP key to verify their Debian packages.
# It should be the same key as https://www.postgresql.org/media/keys/ACCC4CF8.asc
RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8

# Add PostgreSQL's repository. It contains the most recent stable release
#  of PostgreSQL.
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list

# Install ``python-software-properties``, ``software-properties-common`` and PostgreSQL 9.3
#  There are some warnings (in red) that show up during the build. You can hide
#  them by prefixing each apt-get statement with DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y python-software-properties software-properties-common postgresql-9.6 postgresql-client-9.6 postgresql-contrib-9.6 supervisor git
RUN add-apt-repository ppa:openjdk-r/ppa && apt-get update && apt install -y openjdk-11-jdk-headless

USER postgres

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-9.3`` package when it was ``apt-get installed``

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
# Note: here we use ``&&\`` to run commands one after the other - the ``\``
#       allows the RUN command to span multiple lines.
RUN    /etc/init.d/postgresql start &&\
    psql --command "create user mirror_grpc WITH password 'mirror_grpc_pass';" &&\
    psql --command "create user mirror_node with SUPERUSER password 'mirror_node_pass'" &&\
    createdb -O mirror_grpc mirror_node &&\
    psql --command "grant connect on database mirror_node to mirror_grpc;" &&\
    psql --command "alter default privileges in schema public grant select on tables to mirror_grpc;" &&\
    psql --command "grant select on all tables in schema public to mirror_grpc;"


# And add ``listen_addresses`` to ``/etc/postgresql/9.6/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/9.6/main/postgresql.conf
RUN echo "host    all             all             172.17.0.1/16           trust" >> /etc/postgresql/9.6/main/pg_hba.conf

USER root
# --------------------------- PosgreSQL END --------------------------- #


# ---------------------------  Supervisord  --------------------------- #

RUN mkdir -p /var/log/supervisor
COPY supervisord.conf supervisord.conf 

# Copy the Rosetta Executable from the Rosetta Builder stage
WORKDIR /var/rosetta
COPY --from=rosetta-builder /tmp/src/hedera-mirror-rosetta/rosetta-executable .
COPY --from=rosetta-builder /tmp/src/hedera-mirror-rosetta/config/application.yml ./config/application.yml

# Copy the Importer Jar from the Importer Builder stage
WORKDIR /var/importer
COPY --from=java-builder /hedera-mirror-node/hedera-mirror-importer/target .

WORKDIR /

# Expose the ports
EXPOSE 5432 5700

ENTRYPOINT [ "supervisord" ]