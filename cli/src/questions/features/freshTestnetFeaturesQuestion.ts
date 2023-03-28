import {FeaturesQuestion} from "./featuresQuestion";
import {CLIQuestion} from "../question";
import {Answers} from "inquirer";
import {CLIConfig} from "../../config/config";
import {CustomAddressToGiveEGLDQuestion} from "./customAddressToGiveEGLDQuestion";

export class FreshTestnetFeaturesQuestion extends FeaturesQuestion {

    static readonly giveToCustomAddressChoice = 'Give 1,000,000 EGLD to a custom address (otherwise a new one will be generated for you)'
    cliChoices = [FreshTestnetFeaturesQuestion.giveToCustomAddressChoice]

    override async handleAnswer(answers: Answers, config: CLIConfig): Promise<CLIQuestion[]> {
        const questions = await super.handleAnswer(answers, config)

        if (answers.choice.includes(FreshTestnetFeaturesQuestion.giveToCustomAddressChoice)) {
            questions.push(new CustomAddressToGiveEGLDQuestion())
        }

        return questions
    }
}
