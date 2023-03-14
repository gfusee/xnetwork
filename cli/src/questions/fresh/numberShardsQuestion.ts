import {CLIQuestion} from "../question.js";
import {FreshTestnetFeaturesQuestion} from "../features/freshTestnetFeaturesQuestion.js";
import {Answers, Question} from "inquirer";
import {CLIConfig} from "../../config/config.js";

export class NumberShardsQuestion extends CLIQuestion {

    override async getQuestion(): Promise<Question> {
        return {
            type: 'number',
            name: 'numberShards',
            message: 'How many shards do you want to create ? (metachain excluded)',
            default: 1
        }
    }

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        const maxShards = 3

        if (answers.numberShards < 1) {
            const errorMessage = 'Number of shards must be greater than 0'
            console.error(errorMessage)
            return [new NumberShardsQuestion()]
        }

        if (answers.numberShards > maxShards) {
            const errorMessage = `Maximum number of shards is ${maxShards}`
            console.error(errorMessage)
            return [new NumberShardsQuestion()]
        }

        config.numberOfShards = answers.numberShards

        return [new FreshTestnetFeaturesQuestion()]
    }
}
