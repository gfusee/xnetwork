import {CLIChoice, CLIQuestion} from "./question"
import inquirer, {Answers, ListQuestion, Question} from "inquirer";
import {CLIConfig} from "../config/config";
import {NumberShardsQuestion} from "./fresh/numberShardsQuestion";
import {RunnerQuestion} from "./runner/runnerQuestion";
import {removeExistingNetwork} from "../utils/docker/removeExistingNetwork";
import {ContainerState, getNetworkState} from "../utils/docker/getNetworkState";
import {pauseExistingNetwork} from "../utils/docker/pauseExistingNetwork";
import {resumeExistingNetwork} from "../utils/docker/resumeExistingNetwork";
import chalk from "chalk";

export class StartQuestion extends CLIQuestion {

    static readonly removeNetworkChoice = 'Remove existing network'
    static readonly pauseNetworkChoice = 'Pause existing network'
    static readonly resumeNetworkChoice = 'Resume existing network'
    static readonly createNetworkChoice = 'Create a new network'

    override async getQuestion(): Promise<Question> {
        let cliChoices: CLIChoice[] = []
        let cliChoiceMessage = ''

        const state = await getNetworkState()

        if (state.testnetContainerState !== ContainerState.NonExistent) {

            if (state.testnetContainerState === ContainerState.Running) {
                cliChoices.push(StartQuestion.pauseNetworkChoice)
            } else if (state.testnetContainerState === ContainerState.Stopped) {
                cliChoices.push(StartQuestion.resumeNetworkChoice)
            }

            cliChoices.push(StartQuestion.removeNetworkChoice)

            cliChoices.push(new inquirer.Separator())

            const stateChalk = state.testnetContainerState === ContainerState.Running ? chalk.bold.green('running') : chalk.bold.yellow('paused')
            cliChoiceMessage = `A network is ${stateChalk}, what do you want to do ?`
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

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[]> {
        switch (answers.choice) {
            case StartQuestion.removeNetworkChoice:
                await removeExistingNetwork()
                break
            case StartQuestion.pauseNetworkChoice:
                await pauseExistingNetwork()
                break
            case StartQuestion.resumeNetworkChoice:
                await resumeExistingNetwork()
                break
            case StartQuestion.createNetworkChoice:
                return [new NumberShardsQuestion(), new RunnerQuestion()]
        }

        throw 'Unknown choice'
    }

    override async process(config: CLIConfig) {
        if ((await getNetworkState()).testnetContainerState !== ContainerState.NonExistent) {
            await super.process(config)
        } else {
            await (new NumberShardsQuestion()).process(config)
            await (new RunnerQuestion()).process(config)
        }

    }
}
