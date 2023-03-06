FROM --platform=linux/x86_64 ghcr.io/gfusee/mx-testnet-docker:main

COPY wait-for-it.sh .
RUN sudo chmod +x wait-for-it.sh
COPY run.sh .
COPY change-prefs.py change-prefs.py

COPY replace-testnet-toml.py replace-testnet-toml.py

COPY genesis.json genesis.json
COPY create-wallet.py create-wallet.py
COPY change-genesis.py change-genesis.py
COPY economics.toml economics.toml
COPY replace-economics.py replace-economics.py
COPY add-result.py add-result.py
COPY systemSmartContractsConfig.toml systemSmartContractsConfig.toml
COPY replace-system-contracts-config.py replace-system-contracts-config.py

CMD sh run.sh





