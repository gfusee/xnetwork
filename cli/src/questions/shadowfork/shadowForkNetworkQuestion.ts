import {ListQuestion, Question} from "inquirer"
import {CLIQuestion} from "../question.js"
import {ShadowForkFeaturesQuestion} from "../features/shadowForkFeaturesQuestion";

export class ShadowForkTypeQuestion extends CLIQuestion {
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
            onAnswer: () => {},
            nextQuestion: new ShadowForkFeaturesQuestion()
        },
        {
            choice: 'Devnet',
            onAnswer: () => {},
            nextQuestion: new ShadowForkFeaturesQuestion()
        }
    ]
}
