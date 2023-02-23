FROM --platform=linux/x86_64 ghcr.io/gfusee/mx-testnet-docker:main

#COPY mx-chain-db-archive-scripts mx-chain-db-archive-scripts
#RUN mkdir db-downloads
#RUN cd mx-chain-db-archive-scripts && yes | pip3 install -r ./requirements.txt
#RUN cd mx-chain-db-archive-scripts && yes | pip3 install requests
#COPY download-epochs.sh .
#RUN sudo chmod +x download-epochs.sh
#RUN ./download-epochs.sh 0
#RUN ./download-epochs.sh 1
#RUN ./download-epochs.sh 2
#RUN ./download-epochs.sh metachain

COPY mx-chain-tools-go mx-chain-tools-go
RUN cd mx-chain-tools-go/dbMerger/cmd/generalDBMerger && sudo /home/ubuntu/multiversx-sdk/golang/go1.18.4/go/bin/go build

COPY mx-chain-devnet-config mx-chain-devnet-config

COPY merge-db.sh .
RUN sudo chmod +x merge-db.sh

COPY copy-from-extracted.sh .
RUN sudo chmod +x copy-from-extracted.sh
COPY import-db.sh .
RUN sudo chmod +x import-db.sh

COPY wait-for-it.sh .
RUN sudo chmod +x wait-for-it.sh
COPY run.sh .
RUN echo "16"
COPY change-prefs.py change-prefs.py

RUN rm -rf testnet

CMD sh run.sh





