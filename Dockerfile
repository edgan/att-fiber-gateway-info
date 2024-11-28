FROM ubuntu:24.04

COPY docker/att-fiber-gateway-info.sh /usr/local/bin/att-fiber-gateway-info.sh
COPY bin/att-fiber-gateway-info_linux_amd64 /usr/local/bin/att-fiber-gateway-info
COPY default_config.yml /root/.att-fiber-gateway-info.yml

CMD ["att-fiber-gateway-info.sh"]
