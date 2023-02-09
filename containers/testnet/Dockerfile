FROM ubuntu:22.04 AS environment

SHELL ["/bin/bash", "-c"]

RUN apt update -y
RUN apt install sudo -y

RUN adduser --disabled-password --gecos '' ubuntu
RUN adduser ubuntu sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
USER ubuntu
WORKDIR /home/ubuntu

#mxpy setup
RUN sudo apt install build-essential -y
RUN sudo apt install software-properties-common -y
RUN sudo add-apt-repository ppa:deadsnakes/ppa -y
RUN sudo DEBIAN_FRONTEND=noninteractive apt install python3.10 -y
RUN sudo apt install python3.10-venv -y
RUN sudo apt install libncurses5 -y
RUN sudo apt install wget -y
RUN sudo apt install curl -y

RUN sudo wget -O mxpy-up.py https://raw.githubusercontent.com/multiversx/mx-sdk-py-cli/main/mxpy-up.py
RUN python3 mxpy-up.py
ENV PATH="$PATH:/home/ubuntu/multiversx-sdk"

RUN mxpy testnet prerequisites
COPY testnet.toml .
RUN mxpy testnet config

RUN mxpy config set chainID local-testnet
RUN mxpy config set proxy http://localhost:7950

ARG CACHEBUST=1
RUN echo "$CACHEBUST"

COPY wait-for-it.sh .
RUN sudo chmod +x wait-for-it.sh
COPY run.sh .
COPY enable-elastic.py enable-elastic.py

CMD sh run.sh





