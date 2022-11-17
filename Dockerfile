FROM amd64/busybox

LABEL Name="jetbrains_license_servers"
LABEL Version="0.1"
LABEL Description="查找jetbrains license 服务器"

WORKDIR /app

COPY ./bin/jetbrains_license_servers .
COPY ./.env .


CMD ["./jetbrains_license_servers"]
