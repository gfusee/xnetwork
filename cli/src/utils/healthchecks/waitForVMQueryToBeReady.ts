import axios from "axios"
import ora, {Ora} from "ora";

enum NodeStatus {
    NotStarted,
    Starting,
    NotReadyForVmQueries,
    MissingSystemContract,
    Ready
}

export async function waitForVMQueryToBeReady() {
    const spinners: { [key in NodeStatus]: Ora } = {
        [NodeStatus.NotStarted]: ora("Waiting for nodes to respond..."),
        [NodeStatus.Starting]: ora("Waiting for nodes to start committing blocks..."),
        [NodeStatus.NotReadyForVmQueries]: ora("Waiting for nodes to be ready for VM queries..."),
        [NodeStatus.MissingSystemContract]: ora("Waiting for nodes to deploy system smart contracts..."),
        [NodeStatus.Ready]: ora("Nodes are ready for VM queries!")
    }

    const spinnersSuccessMessages = {
        [NodeStatus.NotStarted]: "Nodes are responding",
        [NodeStatus.Starting]: "Nodes are committing blocks",
        [NodeStatus.NotReadyForVmQueries]: "Nodes are ready for VM queries",
        [NodeStatus.MissingSystemContract]: "Nodes have deployed system smart contracts",
        [NodeStatus.Ready]: "Nodes are ready for VM queries"
    }

    let currentStatus = await getNodeStatus()
    let targetStatus = NodeStatus.Ready

    while (currentStatus < targetStatus) {
        succeedPreviousStatusSpinner(currentStatus, spinners, spinnersSuccessMessages)

        if (!spinners[currentStatus].isSpinning) {
            spinners[currentStatus].start()
        }

        await new Promise(resolve => setTimeout(resolve, 1000))
        currentStatus = await getNodeStatus()
    }

    succeedPreviousStatusSpinner(currentStatus, spinners, spinnersSuccessMessages)
}

async function getNodeStatus(): Promise<NodeStatus> {
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

    let queryStatus1: NodeStatus
    let queryStatus2: NodeStatus

    try {
        const queryResponse1 = await axios.post(queryUrl, queryPayload1, { validateStatus: () => true })
        const queryResponse2 = await axios.post(queryUrl, queryPayload2, { validateStatus: () => true })

        if (queryResponse1.status === 200) {
            queryStatus1 = NodeStatus.Ready
        } else {
            queryStatus1 = getNodeStatusFromError(queryResponse1.data.error)
        }

        if (queryResponse2.status === 200) {
            queryStatus2 = NodeStatus.Ready
        } else {
            queryStatus2 = getNodeStatusFromError(queryResponse2.data.error)
        }

        return Math.min(queryStatus1, queryStatus2)
    } catch (e) {
        return NodeStatus.NotStarted
    }
}

function getNodeStatusFromError(error: string): NodeStatus {
    if (error.includes("node is starting")) {
        return NodeStatus.Starting
    }

    if (error.includes("node is not ready yet to process VM Queries")) {
        return NodeStatus.NotReadyForVmQueries
    }

    if (error === "executeQuery:executeQuery:missingsystemsmartcontractonselectedaddress") {
        return NodeStatus.MissingSystemContract
    }

    return NodeStatus.NotStarted
}

function succeedPreviousStatusSpinner(currentStatus: NodeStatus, allSpinners: { [key in NodeStatus] : Ora }, successMessages: { [key in NodeStatus] : string }) {
    const previousStatuses = (Object.values(NodeStatus) as NodeStatus[]).filter(status => status < currentStatus)
    for (const previousStatus of previousStatuses) {
        const spinner = allSpinners[previousStatus]
        if (spinner.isSpinning) {
            spinner.succeed(successMessages[previousStatus])
        }
    }
}

