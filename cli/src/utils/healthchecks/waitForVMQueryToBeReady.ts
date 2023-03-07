import axios from "axios"

export async function waitForVMQueryToBeReady() {
    const url = "http://localhost:7950/vm-values/query"
    const payload = {
        args: [],
        caller: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllst77y4l",
        funcName: "getQueueRegisterNonceAndRewardAddress",
        scAddress: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7"
    }

    while (true) {
        try {
            const response = await axios.post(url, payload, { validateStatus: () => true })
            if (response.status === 200) {
                break
            }
        } catch (error) {
            // do nothing
        }

        await new Promise(resolve => setTimeout(resolve, 1000))
    }
}
