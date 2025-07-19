FROM golang:latest

USER root

RUN apt -y update && apt -y upgrade

RUN apt -y install build-essential pkg-config g++ git cmake yasm

RUN apt install build-essential pkg-config git

RUN apt install -y libc6 libc-bin

RUN apt -y install systemd

RUN apt -y install systemctl

WORKDIR /playground/

COPY docs /playground/

COPY app /playground/

COPY configs/ /playground/configs/

COPY docs /playground/

COPY app.service /etc/systemd/system/

COPY haproxy.cfg /etc/haproxy/haproxy.cfg

RUN chmod +x app

RUN systemctl daemon-reload

RUN ./app migrate

EXPOSE 7500

RUN systemctl start app

# CMD ["systemctl","start","app"]

CMD ["systemctl","start","app"]