import {CLIQuestion} from "../question";
import {Answers, Question} from "inquirer";
import {CLIConfig} from "../../config/config";
import fs from "fs";

export class MxOpsScenesPathQuestion extends CLIQuestion {

    override async getQuestion(): Promise<Question> {
        return {
            type: 'text',
            name: 'mxopsPath',
            message: 'Specify the path where MxOps scenes are located :',
            default: './mxops'
        }
    }

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[]> {
        const answer = answers.mxopsPath

        if (!fs.existsSync(answer)) {
            const errorMessage = 'Path does not exist.'
            console.error(errorMessage)

            return [new MxOpsScenesPathQuestion()]
        }

        config.mxOpsScenesPath = answer

        return []
    }
}
