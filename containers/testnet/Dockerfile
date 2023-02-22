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

COPY extract-epochs.py .
COPY import-db.sh .
RUN sudo chmod +x import-db.sh

COPY wait-for-it.sh .
RUN sudo chmod +x wait-for-it.sh
COPY run.sh .
COPY change-prefs.py change-prefs.py

RUN rm -rf testnet

CMD sh run.sh





