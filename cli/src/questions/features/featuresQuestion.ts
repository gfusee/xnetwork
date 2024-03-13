import {CLIChoice, CLIQuestion} from "../question.js";
import {Answers, CheckboxQuestion, Question} from "inquirer";
import {CLIConfig} from "../../config/config";
import {MxOpsScenesPathQuestion} from "./mxOpsScenesPathQuestion.js";
import {CustomAddressToGiveEGLDQuestion} from "./customAddressToGiveEGLDQuestion.js";

export class FeaturesQuestion extends CLIQuestion {

    private static apiChoice = 'Enable full API (like api.multiversx.com)'
    private static xExplorerChoice = 'Enable xExplorer, requires the full API'
    private static mxOpsChoice = 'Run MxOps scenes at startup'
    static readonly giveToCustomAddressChoice = 'Give 1,000,000 EGLD to a custom address (otherwise a new one will be generated for you)'

    override async getQuestion(): Promise<Question> {
        const question: CheckboxQuestion = {
            type: 'checkbox',
            name: 'choice',
            message: 'Which features do you want to enable ?',
            choices: this.generalFeatures
        }

        return question
    }

    private generalFeatures: CLIChoice[] = [
        FeaturesQuestion.apiChoice,
        FeaturesQuestion.xExplorerChoice,
        FeaturesQuestion.mxOpsChoice,
        FeaturesQuestion.giveToCustomAddressChoice
    ]

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[]> {
        const hasAPI = answers.choice.includes(FeaturesQuestion.apiChoice)
        if (answers.choice.includes(FeaturesQuestion.xExplorerChoice)) {
            if (hasAPI) {
                config.shouldHaveXExplorer = true
            } else {
                const errorMessage = '\n❌ You have to opt for the full API in order to run xExplorer.\n'
                console.log(errorMessage)

                return [new FeaturesQuestion()]
            }
        }

        if (hasAPI) {
            config.shouldHaveElasticSearch = true
            config.shouldHaveMySQL = true
            config.shouldHaveRabbitMQ = true
            config.shouldHaveRedis = true
            config.shouldHaveApi = true
        }

        const questions: CLIQuestion[] = []

        if (answers.choice.includes(FeaturesQuestion.mxOpsChoice)) {
            questions.push(new MxOpsScenesPathQuestion())
        }

        if (answers.choice.includes(FeaturesQuestion.giveToCustomAddressChoice)) {
            questions.push(new CustomAddressToGiveEGLDQuestion())
        }

        return questions
    }
}
