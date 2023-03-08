import {ListQuestion, Question} from "inquirer"
import {CLIQuestion} from "../question.js"
import {ShadowForkFeaturesQuestion} from "../features/shadowForkFeaturesQuestion.js";

export class ShadowForkNetworkQuestion extends CLIQuestion {
    get question(): Question {
        const question: ListQuestion = {
            type: 'list',
            name: 'choice',
            message: 'Which network do you want to fork ?',
            choices: this.cliChoices.map(cliChoice => cliChoice.choice)
        }

        return question
    }

    cliChoices = [
        {
            choice: 'Mainnet',
            nextQuestions: [new ShadowForkFeaturesQuestion()]
        },
        {
            choice: 'Devnet',
            nextQuestions: [new ShadowForkFeaturesQuestion()]
        }
    ]
}
