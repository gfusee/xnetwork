import {execCustomInRepo} from "../exec.js";
import ora from "ora";

export async function pauseExistingNetwork() {
    const removingNetworkSpinner = ora('Pausing network...').start()
    try {
        await execCustomInRepo(`docker-compose exec testnet bash /home/ubuntu/pause.sh`, false)
    } catch (e) {
        console.log(e)
    }
    removingNetworkSpinner.succeed('Paused network successfully')
}
