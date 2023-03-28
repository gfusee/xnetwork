import * as fs from "fs/promises"
import {Constants} from "../../config/constants"
import {CLIConfig} from "../../config/config";

export async function readLatestConfig(): Promise<CLIConfig> {
    const content = await fs.readFile(Constants.CLI_USER_STORAGE_LATEST_CONFIG)

    return JSON.parse(content.toString('utf8')) as CLIConfig
}
