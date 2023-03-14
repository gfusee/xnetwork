import {execCustomInRepo} from "../exec.js";
import ora from "ora";

export async function removeExistingNetwork() {
    const removingNetworkSpinner = ora('Removing the previous network...').start()
    await execCustomInRepo(`docker-compose down`, false)
    removingNetworkSpinner.succeed('Removed the previous network successfully')
}
