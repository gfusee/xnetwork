FROM ubuntu:22.04 AS environment

SHELL ["/bin/bash", "-c"]

RUN apt update -y
RUN apt install sudo -y
RUN apt install curl -y

RUN adduser --disabled-password --gecos '' ubuntu
RUN adduser ubuntu sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
USER ubuntu
WORKDIR /home/ubuntu

#NodeJS setup
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash - && sudo apt install nodejs -y

#mx-api-service setup (nodejs)
RUN sudo apt install git -y
RUN sudo apt install ffmpeg -y
RUN sudo npm install -g @nestjs/cli
RUN git clone https://github.com/multiversx/mx-api-service.git -b main
RUN cd mx-api-service && git checkout 93a9abd0033050e46242836c6cf12e7b396e294c
RUN cd mx-api-service && npm install
RUN rm -rf mx-api-service/config/config.devnet.json
COPY api-config.devnet.yaml mx-api-service/config/config.devnet.yaml
RUN cd mx-api-service && npm run init

#Python3 setup
RUN sudo DEBIAN_FRONTEND=noninteractive apt install python3-pip -y
RUN sudo pip3 install requests

COPY wait-for-it.sh .
RUN sudo chmod +x wait-for-it.sh

COPY run.sh .
RUN sudo chmod +x run.sh
CMD ./run.sh





