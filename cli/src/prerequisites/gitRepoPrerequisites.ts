import {Prerequisites} from "./prerequisites"
import {Constants} from "../config/constants"
import fs from 'fs/promises'
import {execCustom, execCustomInRepo} from "../utils/exec"
import ora from "ora"

export class GitRepoPrerequisites implements Prerequisites {

    async check() {
        if (await this.needToDeleteFiles()) {
            await this.deleteRepo()
        }

        const isRepoAlreadyDownloaded = await this.checkDirectory(Constants.CLI_USER_REPO_PATH)
        const gitCloneSpinner = ora('Checking if necessary datas are downloaded...').start()

        if (!isRepoAlreadyDownloaded) {
            try {
                gitCloneSpinner.info(`Downloading necessary files from ${Constants.REPO_URL}...`)
                await execCustom(`git clone --branch ${Constants.REPO_BRANCH} ${Constants.REPO_URL} ${Constants.CLI_USER_REPO_PATH}`)
                gitCloneSpinner.succeed('Files downloaded')
            } catch (error) {
                gitCloneSpinner.fail(`Error while cloning xnetwork repo: ${Constants.REPO_URL}, branch/tag = ${Constants.REPO_BRANCH}`)
                console.error(error)
                throw error
            }
        } else {
            gitCloneSpinner.succeed('Files already downloaded')
        }
    }

    private async checkDirectory(path: string): Promise<boolean> {
        try {
            const stat = await fs.stat(path)
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

    private async needToDeleteFiles(): Promise<boolean> {
        if (Constants.isNoCacheRequested()) {
            return true
        }

        const gitStatusSpinner = ora('Checking if files are up to date...').start()
        try {
            const gitCurrentTag = (await execCustomInRepo(`git describe --tags --abbrev=0`, false, {cwd: Constants.CLI_USER_REPO_PATH})).stdout.toString()
            const gitCurrentTagVersion = gitCurrentTag.replace('v', '').trim()
            const packageVersion = Constants.getPackageJson().version.trim()
            if (gitCurrentTagVersion !== packageVersion) {
                gitStatusSpinner.fail(`Files are not up to date. Current files version: ${gitCurrentTagVersion} - Package version: ${packageVersion}. Deleting the repo and downloading it again...`)
                return true
            }

            gitStatusSpinner.succeed('Files are up to date')
        } catch (error) {
            if (!(await this.checkDirectory(Constants.CLI_USER_REPO_PATH))) {
                gitStatusSpinner.succeed('Files are not downloaded yet')
                return true
            }
            gitStatusSpinner.fail('Error while checking if repo is up to date')
            throw error
        }

        return false
    }

}
