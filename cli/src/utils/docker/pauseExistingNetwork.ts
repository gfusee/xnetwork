import {execCustomInRepo} from "../exec.js";
import ora from "ora";
import {readLatestConfig} from "../config/readLatestConfig.js"
import {Constants, ContainerInfos, PauseBehavior} from "../../config/constants.js"

export async function pauseExistingNetwork() {
    const latestConfig = await readLatestConfig()

    if (latestConfig.shouldHaveApi) {
        await pauseContainer(Constants.API_CONTAINER)
    }

    if (latestConfig.shouldHaveElasticSearch) {
        await pauseContainer(Constants.ELASTIC_CONTAINER)
    }

    if (latestConfig.shouldHaveMySQL) {
        await pauseContainer(Constants.MYSQL_CONTAINER)
    }

    if (latestConfig.shouldHaveRedis) {
        await pauseContainer(Constants.REDIS_CONTAINER)
    }

    if (latestConfig.shouldHaveRabbitMQ) {
        await pauseContainer(Constants.RABBITMQ_CONTAINER)
    }

    const pausingTestnetSpinner = ora(`Pausing testnet...`).start()
    await pauseContainer(Constants.TESTNET_CONTAINER)
    pausingTestnetSpinner.succeed(`Paused testnet successfully`)
}

async function pauseContainer(containerInfos: ContainerInfos) {
    const command = containerInfos.pauseBehavior === PauseBehavior.PAUSE ? 'docker-compose stop' : 'yes | docker-compose rm -s -v'
    const pausingVerb = containerInfos.pauseBehavior === PauseBehavior.PAUSE ? 'Pausing' : 'Stopping'
    const pausedVerb = containerInfos.pauseBehavior === PauseBehavior.PAUSE ? 'Paused' : 'Stopped'

    const removingNetworkSpinner = ora(`${pausingVerb} ${containerInfos.name}...`).start()
    try {
        await execCustomInRepo(`${command} ${containerInfos.name}`, false)
    } catch (e) {
        console.log(e)
        throw e
    }
    removingNetworkSpinner.succeed(`${pausedVerb} ${containerInfos.name} successfully`)
}
