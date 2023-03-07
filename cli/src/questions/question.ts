import inquirer, {Answers, Question} from "inquirer"
import {CLIConfig} from "../config/config";

export type CLIChoice = {choice: string, nextQuestions: CLIQuestion[] | undefined, onAnswer?: (config: CLIConfig) => void | Promise<void>}

export abstract class CLIQuestion {

    abstract question: Question
    abstract cliChoices: CLIChoice[]

    shouldOverrideActionForChoices: boolean = false

    async process(config: CLIConfig) {
        const response = await inquirer.prompt([this.question])
        const cliChoice = this.cliChoices.find(cliChoice => cliChoice.choice === response.choice)

        let nextQuestions: CLIQuestion[] = []

        if (this.shouldOverrideActionForChoices) {
            const actionResults = await this.overrideActionForAnswers(response, config)

            if (actionResults) {
                nextQuestions.push(...actionResults)
            }
        } else {
            if (!cliChoice) {
                throw new Error(`Unknown choice: ${response.choice}`)
            }

            if (cliChoice.onAnswer) {
                await cliChoice.onAnswer(config)
            }

            if (cliChoice.nextQuestions) {
                nextQuestions.push(...cliChoice.nextQuestions)
            }
        }

        for (const nextQuestion of nextQuestions) {
            await nextQuestion.process(config)
        }
    }

    async overrideActionForAnswers(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        return undefined
    }
}
