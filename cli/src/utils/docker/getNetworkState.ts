import {execCustomInRepo} from "../exec";

type TestnetResult = {
    genesisEgldAddress: string,
    genesisEgldPemPath?: string,
}

type NetworkState = {
    testnetContainerState: ContainerState,
    testnetResult?: TestnetResult
}

export enum ContainerState {
    Running, Stopped, NonExistent
}

export async function getNetworkState(): Promise<NetworkState> {
    const containerState = await getTestnetContainerState()

    try {
        const containerResultsRaw = (await execCustomInRepo("docker-compose exec testnet cat /home/ubuntu/results.json")).stdout.toString()
        const containerResults: TestnetResult = JSON.parse(containerResultsRaw)

        return {
            testnetContainerState: containerState,
            testnetResult: containerResults,
        }
    } catch (e) {
        return {
            testnetContainerState: containerState,
            testnetResult: undefined,
        }
    }
}

async function getTestnetContainerState(): Promise<ContainerState> {
    try {
        const stdout = (await execCustomInRepo('docker-compose ps -a -q testnet')).stdout.toString()

        const containerId = stdout.trim()

        if (!containerId) {
            return ContainerState.NonExistent
        }

        const containerNameRunning = (await execCustomInRepo(`docker ps --filter "id=${containerId}" --filter "status=running" --format '{{.Names}}'`)).stdout.toString().trim()

        if (containerNameRunning) {
            return ContainerState.Running
        } else {
            const containerNameExited = (await execCustomInRepo(`docker ps --filter "id=${containerId}" --filter "status=exited" --format '{{.Names}}'`)).stdout.toString().trim()
            if (containerNameExited) {
                return ContainerState.Stopped
            }

            return ContainerState.NonExistent
        }
    } catch (e) {
        return ContainerState.NonExistent
    }
}
