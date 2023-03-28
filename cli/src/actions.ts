import {CLIConfig, getDefaultConfig} from "./config/config"
import {dontIndent} from "./utils/strings/dontIndent";
import chalk from "chalk";
import {StartQuestion} from "./questions/startQuestion";
import fs from "fs/promises";
import {checkPrerequisites} from "./utils/host/checkPrerequisites";
import {createNetwork} from "./utils/docker/createNetwork";
import {removeExistingNetwork} from "./utils/docker/removeExistingNetwork";
import {readLatestConfig} from "./utils/config/readLatestConfig";

export async function startInteractiveAction() {
    await checkPrerequisites()

    const welcomeMessage = `
    Welcome to ${chalk.bold('xNetwork')} CLI! 🔥
    
    This CLI will help you set up a network with the features you want.
    Let's start by asking you a few questions.
    `

    console.log(dontIndent(welcomeMessage))

    const config = getDefaultConfig()
    await (new StartQuestion()).process(config)
}

export async function createNetworkAction(configPath: string) {
    await checkPrerequisites()

    const configRaw = await fs.readFile(configPath, 'utf-8')
    const config = JSON.parse(configRaw) as CLIConfig

    await createNetwork(config)
}

export async function removeNetworkAction() {
    await checkPrerequisites()

    await removeExistingNetwork()
}

export async function generateConfigAction(outputPath: string) {
    const config = getDefaultConfig()
    const configString = JSON.stringify(config, null, 4)

    await fs.writeFile(outputPath, configString)
}

export async function getLatestConfigAction(): Promise<CLIConfig | undefined> {
    try {
        const latestConfig = await readLatestConfig()
        console.log(JSON.stringify(latestConfig, null, 4))

        return latestConfig
    } catch (e) {
        console.log('No latest config found')
    }

    return undefined
}
