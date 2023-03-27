import {CLIQuestion} from "../question.js";
import {Answers, Question} from "inquirer";
import {CLIConfig} from "../../config/config.js";
import {Address} from "@multiversx/sdk-core/out/index.js";

export class CustomAddressToGiveEGLDQuestion extends CLIQuestion {

    override async getQuestion(): Promise<Question> {
        return {
            type: 'text',
            name: 'address',
            message: 'Which erd address you want to give 1,000,000 EGLD to ?'
        }
    }

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[]> {
        const answer = answers.address

        try {
            Address.fromString(answers.address)
        } catch (e) {
            const errorMessage = 'Address is invalid.'
            console.error(errorMessage)

            return [new CustomAddressToGiveEGLDQuestion()]
        }

        config.initialEGLDAddress = answer

        return []
    }
}
