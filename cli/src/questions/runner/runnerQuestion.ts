import {CLIQuestion} from "../question.js"
import {Answers, ListQuestion, Question} from "inquirer"
import {CLIConfig} from "../../config/config.js"
import {execCustomInRepo, ExecError} from "../../utils/exec.js"
import ora from "ora"
import {waitForVMQueryToBeReady} from "../../utils/healthchecks/waitForVMQueryToBeReady.js";
import {waitForAPIToBeReady} from "../../utils/healthchecks/waitForAPIToBeReady.js";

export class RunnerQuestion extends CLIQuestion {

    private static yesChoice = 'Yes'
    private static noChoice = 'No'

    get question(): Question {
        const question: ListQuestion = {
            type: 'list',
            name: 'choice',
            message: 'All set, do you want to run the network ?',
            choices: this.cliChoices.map(cliChoice => cliChoice.choice)
        }

        return question
    }

    shouldOverrideActionForChoices = true
    cliChoices = [
        {
            choice: RunnerQuestion.yesChoice,
            nextQuestions: undefined
        },
        {
            choice: RunnerQuestion.noChoice,
            nextQuestions: undefined
        }
    ]

    override async overrideActionForAnswers(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        if (answers.choice === RunnerQuestion.yesChoice) {
            await this.run(config)
        }

        return undefined
    }

    private async run(config: CLIConfig) {
        try {
            const removingNetworkSpinner = ora('Removing the previous network...').start()
            await execCustomInRepo(`docker-compose down`, false)
            removingNetworkSpinner.succeed('Removed the previous network successfully')


            if (config.shouldHaveElasticSearch) {
                const startingElasticSearchSpinner = ora('Starting ElasticSearch container...').start()
                await execCustomInRepo(`docker-compose up -d elastic`)
                startingElasticSearchSpinner.succeed('Started ElasticSearch container')
            }

            if (config.shouldHaveMySQL) {
                const startingMySQLSpinner = ora('Starting MySQL container...').start()
                await execCustomInRepo(`docker-compose up -d mysql`)
                startingMySQLSpinner.succeed('Started MySQL container')
            }

            if (config.shouldHaveRedis) {
                const startingRedisSpinner = ora('Starting Redis container...').start()
                await execCustomInRepo(`docker-compose up -d redis`)
                startingRedisSpinner.succeed('Started Redis container')
            }

            if (config.shouldHaveRabbitMQ) {
                const startingRabbitMQSpinner = ora('Starting RabbitMQ container...').start()
                await execCustomInRepo(`docker-compose up -d rabbitmq`)
                startingRabbitMQSpinner.succeed('Started RabbitMQ container')
            }

            if (config.shouldHaveApi) {
                const startingApiSpinner = ora('Starting API container...').start()
                await execCustomInRepo(`docker-compose up -d api`)
                startingApiSpinner.succeed('Started API container')
            }

            const startingNetworkSpinner = ora('Starting network...').start()
            await execCustomInRepo(`docker-compose up -d testnet`, false, {
                env: {
                    ...process.env,
                    "MX_LT_NUM_SHARDS": config.numberOfShards.toString(),
                    "MX_LT_ELASTIC_ENABLED": config.shouldHaveElasticSearch.toString(),
                    "MX_LT_CUSTOM_EGLD_ADDRESS": config.initialEGLDAddress ?? ""
                }
            })
            startingNetworkSpinner.succeed('Started network successfully')

            const networkReadySpinner = ora('Waiting for network to be ready...').start()
            await waitForVMQueryToBeReady()
            networkReadySpinner.succeed('Network is ready')

            if (config.shouldHaveApi) {
                const startingApiHealthCheckSpinner = ora('Waiting for API to be ready').start()
                await waitForAPIToBeReady()
                startingApiHealthCheckSpinner.succeed('API is ready')
            }
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
}
