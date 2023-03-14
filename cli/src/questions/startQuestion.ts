import {CLIChoice, CLIQuestion} from "./question.js"
import inquirer, {Answers, ListQuestion, Question} from "inquirer";
import {CLIConfig} from "../config/config.js";
import {NumberShardsQuestion} from "./fresh/numberShardsQuestion.js";
import {execCustomInRepo} from "../utils/exec.js";
import {RunnerQuestion} from "./runner/runnerQuestion.js";
import {removeExistingNetwork} from "../utils/docker/removeExistingNetwork.js";

export class StartQuestion extends CLIQuestion {

    static readonly removeNetworkChoice = 'Remove existing network'
    static readonly createNetworkChoice = 'Create a new network'

    override async getQuestion(): Promise<Question> {
        let cliChoices: CLIChoice[] = []
        let cliChoiceMessage = ''

        if (await this.isNetworkRunning()) {
            cliChoices.push(StartQuestion.removeNetworkChoice)
            cliChoices.push(new inquirer.Separator())
            cliChoiceMessage = 'A network is already running, what do you want to do ?'
        }

        cliChoices.push(StartQuestion.createNetworkChoice)

        const question: ListQuestion = {
            type: 'list',
            name: 'choice',
            message: cliChoiceMessage,
            choices: cliChoices
        }

        return question
    }

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        if (answers.choice === StartQuestion.removeNetworkChoice) {
            await removeExistingNetwork()
        } else if (answers.choice === StartQuestion.createNetworkChoice) {
            return [new NumberShardsQuestion(), new RunnerQuestion()]
        }
    }

    private async isNetworkRunning(): Promise<boolean> {
        try {
            await execCustomInRepo(`docker-compose ps | grep -q " Up "`, false)
            return true
        } catch (e) {
            return false
        }
    }

    override async process(config: CLIConfig) {
        if (await this.isNetworkRunning()) {
            await super.process(config)
        } else {
            await (new NumberShardsQuestion()).process(config)
            await (new RunnerQuestion()).process(config)
        }

    }
}
