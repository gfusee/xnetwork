import {FeaturesQuestion} from "./featuresQuestion.js";
import {CLIQuestion} from "../question.js";
import {Answers} from "inquirer";
import {CLIConfig} from "../../config/config.js";
import {CustomAddressToGiveEGLDQuestion} from "./customAddressToGiveEGLDQuestion.js";

export class FreshTestnetFeaturesQuestion extends FeaturesQuestion {

    static readonly giveToCustomAddressChoice = 'Give 1,000,000 EGLD to a custom address'
    cliChoices = [
        {
            choice: FreshTestnetFeaturesQuestion.giveToCustomAddressChoice,
            nextQuestions: undefined
        },
    ]

    override async overrideActionForAnswers(answers: Answers, config: CLIConfig): Promise<CLIQuestion[] | undefined> {
        await super.overrideActionForAnswers(answers, config)

        if (answers.choice.includes(FreshTestnetFeaturesQuestion.giveToCustomAddressChoice)) {
            return [new CustomAddressToGiveEGLDQuestion()]
        }

        return undefined
    }
}
