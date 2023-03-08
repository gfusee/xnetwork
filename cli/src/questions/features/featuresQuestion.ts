import {CLIChoice, CLIQuestion} from "../question.js";
import {Answers, CheckboxQuestion, Question} from "inquirer";
import {CLIConfig} from "../../config/config";

export abstract class FeaturesQuestion extends CLIQuestion {

    private static apiChoice = 'Enable full API (like api.multiversx.com)'

    shouldOverrideActionForChoices: boolean = true

    get question(): Question {
        const question: CheckboxQuestion = {
            type: 'checkbox',
            name: 'choice',
            message: 'Which features do you want to enable ?',
            choices: this.generalFeatures.concat(this.cliChoices).map(cliChoice => cliChoice.choice)
        }

        return question
    }

    private generalFeatures: CLIChoice[] = [
        {
            choice: FeaturesQuestion.apiChoice,
            nextQuestions: undefined
        }
    ]

    abstract cliChoices: CLIChoice[]

    override async overrideActionForAnswers(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        if (answers.choice.includes(FeaturesQuestion.apiChoice)) {
            config.shouldHaveElasticSearch = true
            config.shouldHaveMySQL = true
            config.shouldHaveRabbitMQ = true
            config.shouldHaveRedis = true
            config.shouldHaveApi = true
        }

        return undefined
    }
}
