import {execCustomInRepo} from "../exec";

export async function getMxOpsScenarioDatas(): Promise<any> {
    const mxopsXNetworkValuesRaw = (await execCustomInRepo(`docker-compose exec testnet python3 -m mxops data get -n LOCAL -s xnetwork`)).stdout.toString()
    const searchString = 'ABSOLUTELY NO WARRANTY\n'
    const mxopsXNetworkValues = mxopsXNetworkValuesRaw.substring(mxopsXNetworkValuesRaw.lastIndexOf(searchString) + searchString.length).trim()
    return JSON.parse(mxopsXNetworkValues)
}
