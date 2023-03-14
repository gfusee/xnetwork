import inquirer, {Answers, Question} from "inquirer"
import {CLIConfig} from "../config/config";

export type CLIChoice = string | inquirer.Separator

export abstract class CLIQuestion {

    abstract getQuestion(): Promise<Question>

    async process(config: CLIConfig) {
        const response = await inquirer.prompt([await this.getQuestion()])

        const nextQuestions: CLIQuestion[] = []
        const actionResults = await this.handleAnswer(response, config)
        if (actionResults) {
            nextQuestions.push(...actionResults)
        }

        for (const nextQuestion of nextQuestions) {
            await nextQuestion.process(config)
        }
    }

    async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        return undefined
    }
}
