import {Prerequisites} from "./prerequisites"
import {Constants} from "../config/constants.js"
import fs from 'fs/promises'
import {execCustom} from "../utils/exec.js"
import ora from "ora"

export class GitRepoPrerequisites implements Prerequisites {

    constructor(private noCache: boolean) {}

    async check() {
        if (this.noCache) {
            await this.deleteRepo()
        }

        const isRepoAlreadyDownloaded = await this.checkDirectory(Constants.CLI_USER_REPO_PATH)
        const gitCloneSpinner = ora('Checking if necessary datas are downloaded...').start()

        if (!isRepoAlreadyDownloaded) {
            try {
                gitCloneSpinner.info(`Downloading necessary files from ${Constants.REPO_URL}...`)
                await execCustom(`git clone ${Constants.REPO_URL} ${Constants.CLI_USER_REPO_PATH}`)
                gitCloneSpinner.succeed('Files downloaded')
            } catch (error) {
                gitCloneSpinner.fail(`Error while cloning xnetwork repo: ${Constants.REPO_URL}`)
                throw error
            }
        } else {
            gitCloneSpinner.succeed('Files already downloaded')
        }
    }

    private async checkDirectory(path: string): Promise<boolean> {
        try {
            const stat = await fs.stat(Constants.CLI_USER_REPO_PATH)
            return stat.isDirectory()
        } catch (error) {
            return false
        }
    }

    private async deleteRepo() {
        const deleteSpinner = ora('Deleting existing files...').start()
        try {
            await fs.rm(Constants.CLI_USER_REPO_PATH, {recursive: true, force: true})
            deleteSpinner.succeed('Files deleted')
        } catch (error) {
            deleteSpinner.fail('Error while deleting files')
            throw error
        }
    }

}
