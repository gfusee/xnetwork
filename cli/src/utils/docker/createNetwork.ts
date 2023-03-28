import {CLIConfig} from "../../config/config";
import {removeExistingNetwork} from "./removeExistingNetwork";
import ora from "ora";
import {upContainer} from "./upContainer";
import {Constants} from "../../config/constants";
import {saveLatestConfig} from "../config/saveLatestConfig";
import {waitForVMQueryToBeReady} from "../healthchecks/waitForVMQueryToBeReady";
import {waitForAPIToBeReady} from "../healthchecks/waitForAPIToBeReady";
import {execCustomInRepo, ExecError} from "../exec";
import {ResultLogger} from "../../result/resultLogger";

export async function createNetwork(config: CLIConfig) {
    try {
        await removeExistingNetwork()

        if (config.shouldHaveElasticSearch) {
            const startingElasticSearchSpinner = ora('Starting ElasticSearch container...').start()
            await upContainer(Constants.ELASTIC_CONTAINER.name)
            startingElasticSearchSpinner.succeed('Started ElasticSearch container')
        }

        if (config.shouldHaveMySQL) {
            const startingMySQLSpinner = ora('Starting MySQL container...').start()
            await upContainer(Constants.MYSQL_CONTAINER.name)
            startingMySQLSpinner.succeed('Started MySQL container')
        }

        if (config.shouldHaveRedis) {
            const startingRedisSpinner = ora('Starting Redis container...').start()
            await upContainer(Constants.REDIS_CONTAINER.name)
            startingRedisSpinner.succeed('Started Redis container')
        }

        if (config.shouldHaveRabbitMQ) {
            const startingRabbitMQSpinner = ora('Starting RabbitMQ container...').start()
            await upContainer(Constants.RABBITMQ_CONTAINER.name)
            startingRabbitMQSpinner.succeed('Started RabbitMQ container')
        }

        const startingNetworkSpinner = ora('Starting network...').start()
        await upContainer(Constants.TESTNET_CONTAINER.name, {
            ...process.env,
            "MX_LT_NUM_SHARDS": config.numberOfShards.toString(),
            "MX_LT_ELASTIC_ENABLED": config.shouldHaveElasticSearch.toString(),
            "MX_LT_CUSTOM_EGLD_ADDRESS": config.initialEGLDAddress ?? ""
        })
        startingNetworkSpinner.succeed('Started network successfully')

        await saveLatestConfig(config)

        await waitForVMQueryToBeReady()

        if (config.shouldHaveApi) {
            const startingApiSpinner = ora('Starting API container...').start()
            await upContainer(Constants.API_CONTAINER.name)
            startingApiSpinner.succeed('Started API container')

            const startingApiHealthCheckSpinner = ora('Waiting for API to be ready').start()
            await waitForAPIToBeReady()
            startingApiHealthCheckSpinner.succeed('API is ready')
        }

        if (config.mxOpsScenesPath) {
            const copyingScenesSpinner = ora('Copying mxops scenes...').start()
            await execCustomInRepo(`docker-compose cp ${config.mxOpsScenesPath} testnet:/home/ubuntu/mxops`)
            copyingScenesSpinner.succeed('Copied mxops scenes')

            const runningScenesSpinner = ora('Running mxops scenes...').start()
            await execCustomInRepo(`docker-compose exec testnet python3 run_mxops.py`)
            runningScenesSpinner.succeed('Ran mxops scenes')
        }

        const resultLogger = new ResultLogger()
        await resultLogger.printResults(config)
    } catch (e) {
        try {
            const error = e as ExecError
            console.log("Error while running network...")
            console.log("Command : ", error.error.cmd)
            console.log("Exit code : ", error.error.code)
            console.log("Message : ", error.stderr)
        } catch (e2) {
            console.log("Error while running network...")
            console.log(e)
        }

        throw e
    }
}
