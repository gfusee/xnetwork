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

#Python3 setup
RUN sudo DEBIAN_FRONTEND=noninteractive apt install python3-pip -y
RUN sudo pip3 install requests

#NodeJS setup
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash - && sudo apt install nodejs -y
RUN sudo npm i -g yarn

#mx-explorer-dapp setup (nodejs)
RUN sudo apt install git -y
RUN git clone https://github.com/multiversx/mx-explorer-dapp.git -b main
RUN cd mx-explorer-dapp && git checkout 6b9560d83cfa9b2197f1efde651450c9d25508e4
RUN cd mx-explorer-dapp && yarn
COPY config.localnet.ts mx-explorer-dapp/src/config/index.ts
COPY vite.config.ts mx-explorer-dapp/vite.config.ts

COPY wait-for-it.sh .
RUN sudo chmod +x wait-for-it.sh

COPY run.sh .
RUN sudo chmod +x run.sh
CMD ./run.sh
