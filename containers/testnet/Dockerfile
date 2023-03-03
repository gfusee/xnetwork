FROM --platform=linux/x86_64 ghcr.io/gfusee/mx-testnet-docker:main

COPY wait-for-it.sh .
RUN sudo chmod +x wait-for-it.sh
COPY run.sh .
COPY change-prefs.py change-prefs.py

COPY replace-testnet-toml.py replace-testnet-toml.py

CMD sh run.sh





