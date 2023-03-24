import * as fs from "fs/promises"
import { existsSync } from "fs"
import {CLIConfig} from "../../config/config.js"
import {Constants} from "../../config/constants.js"
import ora from "ora";

// Write the latest config to the file system, creating it if it doesn't exist

export async function saveLatestConfig(config: CLIConfig) {
    const savingConfigSpinner = ora('Saving config...').start()
    if (!existsSync(Constants.CLI_USER_STORAGE_PATH)) {
        await fs.mkdir(Constants.CLI_USER_STORAGE_PATH)
    }
    await fs.writeFile(Constants.CLI_USER_STORAGE_LATEST_CONFIG, JSON.stringify(config, null, 4))
    savingConfigSpinner.succeed('Saved config successfully')
}
