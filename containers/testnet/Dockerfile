FROM --platform=amd64 ghcr.io/gfusee/mx-testnet-docker:main@sha256:710e39d017fe6be27b3c9a1b61a06e9b168e157a263f19ea89e02b8c0924a2d3

STOPSIGNAL SIGTERM

RUN pip3 install mxops==1.0.0

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
COPY handle_signal.sh handle_signal.sh
COPY mxops-init mxops-init
COPY run_mxops.py run_mxops.py

RUN sudo chmod +x pause.sh
RUN sudo chmod +x handle_signal.sh
RUN sudo chmod +x run.sh

CMD ./run.sh





