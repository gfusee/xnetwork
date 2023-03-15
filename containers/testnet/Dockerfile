FROM --platform=linux/x86_64 ghcr.io/gfusee/mx-testnet-docker:main

RUN sudo apt install lsof -y

COPY wait-for-it.sh .
RUN sudo chmod +x wait-for-it.sh
COPY run.sh .
COPY change_prefs.py change_prefs.py

COPY replace_testnet_toml.py replace_testnet_toml.py

COPY genesis.json genesis.json
COPY create_wallet.py create_wallet.py
COPY change_genesis.py change_genesis.py
COPY economics.toml economics.toml
COPY replace_economics.py replace_economics.py
COPY add_result.py add_result.py
COPY read_result.py read_result.py
COPY systemSmartContractsConfig.toml systemSmartContractsConfig.toml
COPY replace_system_contracts_config.py replace_system_contracts_config.py
COPY pause.sh pause.sh

RUN sudo chmod +x pause.sh

CMD sh run.sh





