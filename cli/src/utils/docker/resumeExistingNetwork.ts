import {execCustomInRepo} from "../exec";
import ora from "ora";
import {readLatestConfig} from "../config/readLatestConfig";
import {upContainer} from "./upContainer";
import {Constants, ContainerInfos, PauseBehavior} from "../../config/constants";
import {waitForVMQueryToBeReady} from "../healthchecks/waitForVMQueryToBeReady";

export async function resumeExistingNetwork() {
    const latestConfig = await readLatestConfig()

    if (latestConfig.shouldHaveElasticSearch) {
        await resumeContainer(Constants.ELASTIC_CONTAINER)
    }

    if (latestConfig.shouldHaveMySQL) {
        await resumeContainer(Constants.MYSQL_CONTAINER)
    }

    if (latestConfig.shouldHaveRedis) {
        await resumeContainer(Constants.REDIS_CONTAINER)
    }

    if (latestConfig.shouldHaveRabbitMQ) {
        await resumeContainer(Constants.RABBITMQ_CONTAINER)
    }

    await resumeContainer(Constants.TESTNET_CONTAINER)

    await waitForVMQueryToBeReady()

    if (latestConfig.shouldHaveApi) {
        await resumeContainer(Constants.API_CONTAINER)
    }
}

async function resumeContainer(containerInfos: ContainerInfos) {
    const removingNetworkSpinner = ora(`Resuming ${containerInfos.name}...`).start()

    if (containerInfos.pauseBehavior === PauseBehavior.PAUSE) {
        await execCustomInRepo(`docker-compose start ${containerInfos.name}`, false)
    } else {
        await upContainer(containerInfos.name, undefined, false)
    }

    removingNetworkSpinner.succeed(`Resumed ${containerInfos.name} successfully`)
}
