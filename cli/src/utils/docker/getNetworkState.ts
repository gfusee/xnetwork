import {execCustomInRepo} from "../exec.js";

type LocalnetResult = {
    genesisEgldAddress: string,
    genesisEgldPemPath?: string,
}

type NetworkState = {
    localnetContainerState: ContainerState,
    localnetResult?: LocalnetResult
}

export enum ContainerState {
    Running, Stopped, NonExistent
}

export async function getNetworkState(): Promise<NetworkState> {
    const containerState = await getLocalnetContainerState()

    try {
        const containerResultsRaw = (await execCustomInRepo("docker compose exec localnet cat /home/ubuntu/results.json")).stdout.toString()
        const containerResults: LocalnetResult = JSON.parse(containerResultsRaw)

        return {
            localnetContainerState: containerState,
            localnetResult: containerResults,
        }
    } catch (e) {
        return {
            localnetContainerState: containerState,
            localnetResult: undefined,
        }
    }
}

async function getLocalnetContainerState(): Promise<ContainerState> {
    try {
        const stdout = (await execCustomInRepo('docker compose ps -a -q localnet')).stdout.toString()

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
