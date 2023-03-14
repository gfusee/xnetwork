import {FeaturesQuestion} from "./featuresQuestion.js";
import {CLIQuestion} from "../question.js";
import {Answers} from "inquirer";
import {CLIConfig} from "../../config/config.js";
import {CustomAddressToGiveEGLDQuestion} from "./customAddressToGiveEGLDQuestion.js";

export class FreshTestnetFeaturesQuestion extends FeaturesQuestion {

    static readonly giveToCustomAddressChoice = 'Give 1,000,000 EGLD to a custom address (otherwise a new one will be generated for you)'
    cliChoices = [FreshTestnetFeaturesQuestion.giveToCustomAddressChoice]

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        await super.handleAnswer(answers, config)

        if (answers.choice.includes(FreshTestnetFeaturesQuestion.giveToCustomAddressChoice)) {
            return [new CustomAddressToGiveEGLDQuestion()]
        }

        return undefined
    }
}
