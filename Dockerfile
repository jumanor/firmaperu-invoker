# Compilacion
FROM golang:1.19 AS builder

WORKDIR /opt
COPY . .

RUN go mod tidy
RUN go build -o main main.go

# Usa la imagen base httpd:2.4
FROM httpd:2.4

COPY --from=builder /opt/main /opt/main
RUN chmod +x /opt/main
COPY ./public/ /opt/public
COPY ./config.properties.example /opt/config.properties.example
COPY ./example/ /usr/local/apache2/htdocs
COPY entrypoint.sh /usr/local/bin
RUN chmod +x /usr/local/bin/entrypoint.sh

ENV CLIENT_ID=""
ENV CLIENT_SECRET=""


RUN sed -i '/^Listen 80/a Listen 5050' /usr/local/apache2/conf/httpd.conf

# Servidor Web
EXPOSE 5050
# Servidor firmaperu-invoker
EXPOSE 9091

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Actualiza el índice de paquetes y luego instala las dependencias necesarias
RUN apt-get update && \
    apt-get install -y p7zip-full && \
    rm -rf /var/lib/apt/lists/*




