import {CLIChoice, CLIQuestion} from "../question.js";
import {Answers, CheckboxQuestion, Question} from "inquirer";
import {CLIConfig} from "../../config/config";

export abstract class FeaturesQuestion extends CLIQuestion {

    private static apiChoice = 'Enable full API (like api.multiversx.com)'

    override async getQuestion(): Promise<Question> {
        const question: CheckboxQuestion = {
            type: 'checkbox',
            name: 'choice',
            message: 'Which features do you want to enable ?',
            choices: this.generalFeatures.concat(this.cliChoices)
        }

        return question
    }

    private generalFeatures: CLIChoice[] = [FeaturesQuestion.apiChoice]

    abstract cliChoices: CLIChoice[]

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
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
