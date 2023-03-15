import {execCustomInRepo} from "../exec.js";
import ora from "ora";

export async function resumeExistingNetwork() {
    const removingNetworkSpinner = ora('Resuming the network...').start()
    await execCustomInRepo(`docker-compose up -d testnet`, false)
    removingNetworkSpinner.succeed('Resumed network successfully')
}
