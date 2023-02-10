# MX full local testnet environment

This project provides an easy way to run a local testnet of MultiversX and its associated API, with no configuration required.
While you can run a local testnet of MultiversX with the [MultiversX CLI](
), this project provides a more complete environment, with a proxy and a full API.
## Features
- Run your own testnet of MultiversX, locally
- Run the associated API, the same as https://api.multiversx.com

## Requirements
- [Docker](https://docs.docker.com/get-started/) and [Docker Compose](https://docs.docker.com/compose/gettingstarted/).

## Usage

1. Clone this repository

    ```bash
    git clone https://github.com/gfusee/mx-full-testnet-docker.git
    ```
   
2. Change into the project directory

    ```bash
    cd mx-full-testnet-docker
    ```
   
3. Start the stack with the following command:

    ```bash
    docker-compose up -d
    ```
   
4. Wait for the stack to start. This can take a few minutes. You can check the status of the stack with the following command (Ctrl+C to exit):
    ```bash
    docker-compose logs api -t -f
    ```
   
Voilà ! After the setup, you can access the API at http://localhost:3001.
The proxy is also exposed at http://localhost:7950.
   
# Stopping the stack

To stop the stack, run the following command:

```bash
docker-compose down
```

Do not stop the stack with ```docker-compose stop```, as this will not remove the volumes, it'll make the blockchain not sync with other services and leads to errors.
# Contributing

If you find any issues or would like to contribute to this project, feel free to open a pull request or an issue on GitHub.
