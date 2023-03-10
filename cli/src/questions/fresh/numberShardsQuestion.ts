import {CLIQuestion} from "../question.js";
import {FreshTestnetFeaturesQuestion} from "../features/freshTestnetFeaturesQuestion.js";
import {Answers} from "inquirer";
import {CLIConfig} from "../../config/config.js";

export class NumberShardsQuestion extends CLIQuestion {
    shouldOverrideActionForChoices: boolean = true
    question = {
        type: 'number',
        name: 'numberShards',
        message: 'How many shards do you want to create ? (metachain excluded)',
        default: 1
    }

    cliChoices = []

    override async overrideActionForAnswers(choices: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        const maxShards = 3

        if (choices.numberShards < 1) {
            const errorMessage = 'Number of shards must be greater than 0'
            console.error(errorMessage)
            return [new NumberShardsQuestion()]
        }

        if (choices.numberShards > maxShards) {
            const errorMessage = `Maximum number of shards is ${maxShards}`
            console.error(errorMessage)
            return [new NumberShardsQuestion()]
        }

        config.numberOfShards = choices.numberShards

        return [new FreshTestnetFeaturesQuestion()]
    }
}
