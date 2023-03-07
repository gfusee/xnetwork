import {CLIQuestion} from "../question.js";
import {Answers} from "inquirer";
import {CLIConfig} from "../../config/config.js";
import {Address} from "@multiversx/sdk-core/out/index.js";

export class CustomAddressToGiveEGLDQuestion extends CLIQuestion {
    shouldOverrideActionForChoices: boolean = true
    question = {
        type: 'text',
        name: 'address',
        message: 'Which erd address you want to give 1,000,000 EGLD to ?'
    }

    cliChoices = []

    override async overrideActionForAnswers(choices: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        const answer = choices.address

        try {
            Address.fromString(choices.address)
        } catch (e) {
            const errorMessage = 'Address is invalid.'
            console.error(errorMessage)

            return [new CustomAddressToGiveEGLDQuestion()]
        }

        config.initialEGLDAddress = answer

        return undefined
    }
}
