import {CLIConfig} from "../config/config.js";
import chalk from "chalk";
import {execCustom} from "../utils/exec.js";
import {dontIndent} from "../utils/strings/dontIndent.js";

type Result = {
    genesisEgldAddress: string,
    genesisEgldPemPath?: string,
}

export class ResultLogger {
    async printResults(config: CLIConfig) {

        const containerResultsRaw = (await execCustom("docker-compose exec testnet cat /home/ubuntu/results.json")).stdout.toString()
        const containerResults: Result = JSON.parse(containerResultsRaw)

        let resultString = `${chalk.bold.green("Local testnet successfully started !")}`

        resultString += `
        
        Proxy/Gateway URL: ${chalk.blue("http://localhost:7950")}
        ChainID: ${chalk.blue("local-testnet")}
        `

        if (config.shouldHaveElasticSearch) {
            resultString += `
            ElasticSearch URL: ${chalk.blue("http://localhost:9200")}
            `
        }

        if (config.shouldHaveApi) {
            resultString += `
            API URL: ${chalk.blue("http://localhost:3001")}
            `
        }

        if (containerResults.genesisEgldPemPath) {
            const addressPrivateKey = (await execCustom(`docker-compose exec testnet cat ${containerResults.genesisEgldPemPath}`)).stdout.toString()
            resultString += `
            An address with 1,000,000 EGLD was generated for you. Here are the details:
            
            ${chalk.bold.red("Here is the private key. Keep it safe and don't use it in another place than the local testnet!")}
            
            ${addressPrivateKey}
            
            ${chalk.bold.red("End of the private key. You can copy/paste the content into a .pem file and use it to do transactions on the local testnet.")}
            `
        } else {
            resultString += `
            You choose to give 1,000,000 EGLD to a custom address instead of generating a new one. As a reminder, here is the address you specified: ${chalk.blue(containerResults.genesisEgldAddress)}
            `
        }

        console.log(dontIndent(resultString))

    }

}
