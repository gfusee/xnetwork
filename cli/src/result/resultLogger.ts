import {CLIConfig} from "../config/config.js";
import chalk from "chalk";
import {execCustomInRepo} from "../utils/exec.js";
import {dontIndent} from "../utils/strings/dontIndent.js";
import {getNetworkState} from "../utils/docker/getNetworkState.js";

export class ResultLogger {
    async printResults(config: CLIConfig) {

        const state = await getNetworkState()
        const containerResults = state.localnetResult

        if (!containerResults) {
            console.log(dontIndent(chalk.bold.red("Something went wrong while checking the status of the localnet.")))
            return
        }

        let firstPartResultString = `${chalk.bold.green("Localnet successfully started !")}`

        firstPartResultString += `
        
        Proxy/Gateway URL: ${chalk.blue("http://localhost:7950")}
        ChainID: ${chalk.blue("localnet")}
        `

        if (config.shouldHaveElasticSearch) {
            firstPartResultString += `
            ElasticSearch URL: ${chalk.blue("http://localhost:9200")}
            `
        }

        if (config.shouldHaveApi) {
            firstPartResultString += `
            API URL: ${chalk.blue("http://localhost:3001")}
            `
        }

        if (containerResults.genesisEgldPemPath) {
            const addressPrivateKey = (await execCustomInRepo(`docker-compose exec localnet cat ${containerResults.genesisEgldPemPath}`)).stdout.toString()
            firstPartResultString += `
            An address with 1,000,000 EGLD was generated for you. Here are the details:
            
            ${chalk.bold.red("Here is the private key. Keep it safe and don't use it in another place than the localnet!")}
            
            ${addressPrivateKey}
            
            ${chalk.bold.red("End of the private key. You can copy/paste the content into a .pem file and use it to do transactions on the localnet.")}
            `
        } else {
            firstPartResultString += `
            You choose to give 1,000,000 EGLD to a custom address instead of generating a new one. As a reminder, here is the address you specified: ${chalk.blue(containerResults.genesisEgldAddress)}
            `
        }

        let mxOpsDisplayString = ''

        if (config.mxOpsScenesPath) {
            const mxopsXNetworkValuesRaw = (await execCustomInRepo(`docker-compose exec localnet python3 -m mxops data get -n LOCAL -s xnetwork`)).stdout.toString()
            const searchString = 'ABSOLUTELY NO WARRANTY\n'
            const mxopsXNetworkValues = mxopsXNetworkValuesRaw.substring(mxopsXNetworkValuesRaw.lastIndexOf(searchString) + searchString.length).trim()
            const mxopsXNetworkValuesObject = JSON.parse(mxopsXNetworkValues)

            mxOpsDisplayString = dontIndent(
                `
            ------------------------------
            ${chalk.bold.red("You specified a path to mxops scenes. They have been processed under a scenario called 'xnetwork'. Here are result values of this scenario:")}
            `
            )

            mxOpsDisplayString += '\n'
            mxOpsDisplayString += JSON.stringify(mxopsXNetworkValuesObject, null, 4).trim()

            mxOpsDisplayString += dontIndent(
                `
            ------------------------------
            `
            )
        }

        console.log(dontIndent(firstPartResultString))

        if (mxOpsDisplayString.length > 0) {
            console.log(mxOpsDisplayString)
        }

    }

}
