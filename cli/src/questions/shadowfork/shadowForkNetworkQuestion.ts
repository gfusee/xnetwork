import {Answers, ListQuestion, Question} from "inquirer"
import {CLIQuestion} from "../question"
import {ShadowForkFeaturesQuestion} from "../features/shadowForkFeaturesQuestion";

export class ShadowForkNetworkQuestion extends CLIQuestion {
    static readonly mainnetChoice = 'Mainnet'
    static readonly devnetChoice = 'Devnet'

    override async getQuestion(): Promise<Question> {
        const listQuestion: ListQuestion = {
            type: 'list',
            name: 'choice',
            message: 'Which network do you want to fork ?',
            choices: [ShadowForkNetworkQuestion.mainnetChoice, ShadowForkNetworkQuestion.devnetChoice]
        }

        return listQuestion
    }

    override async handleAnswer(answers: Answers): Promise<CLIQuestion[]> {
        return [new ShadowForkFeaturesQuestion()]
    }
}
