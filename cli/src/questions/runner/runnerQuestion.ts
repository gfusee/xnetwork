import {CLIQuestion} from "../question"
import {Answers, ListQuestion, Question} from "inquirer"
import {CLIConfig} from "../../config/config"
import {createNetwork} from "../../utils/docker/createNetwork";

export class RunnerQuestion extends CLIQuestion {

    private static yesChoice = 'Yes'
    private static noChoice = 'No'

    override async getQuestion(): Promise<Question> {
        const listQuestion: ListQuestion = {
            type: 'list',
            name: 'choice',
            message: 'All set, do you want to run the network ?',
            choices: [RunnerQuestion.yesChoice, RunnerQuestion.noChoice]
        }

        return listQuestion
    }

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[]> {
        if (answers.choice === RunnerQuestion.yesChoice) {
            await createNetwork(config)
        }

        return []
    }
}
