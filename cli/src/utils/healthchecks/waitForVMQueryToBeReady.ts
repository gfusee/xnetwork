import axios from "axios"

export async function waitForVMQueryToBeReady() {
    const queryUrl = "http://localhost:7950/vm-values/query"
    const queryPayload1 = {
        args: [],
        caller: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllst77y4l",
        funcName: "getQueueRegisterNonceAndRewardAddress",
        scAddress: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7"
    }
    const queryPayload2 = {
        scAddress: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqylllslmq6y6",
        funcName: "getAllContractAddresses",
        args: []
    }

    const economicsUrl = "http://localhost:7950/network/economics"

    while (true) {
        try {
            const queryResponse1 = await axios.post(queryUrl, queryPayload1, { validateStatus: () => true })
            const queryResponse2 = await axios.post(queryUrl, queryPayload2, { validateStatus: () => true })
            const economicsResponse = await axios.get(economicsUrl, { validateStatus: () => true })
            if (queryResponse1.status === 200 && queryResponse2.status === 200 && economicsResponse.status === 200) {
                break
            }
        } catch (error) {
            // do nothing
        }

        await new Promise(resolve => setTimeout(resolve, 1000))
    }
}
