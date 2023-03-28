#!/usr/bin/env node

import {CLIConfig, getDefaultConfig} from "./config/config.js"
import {dontIndent} from "./utils/strings/dontIndent.js";
import chalk from "chalk";
import { program } from 'commander'
import {StartQuestion} from "./questions/startQuestion.js";
import fs from "fs/promises";
import {checkPrerequisites} from "./utils/host/checkPrerequisites.js";
import {createNetwork} from "./utils/docker/createNetwork.js";
import {removeExistingNetwork} from "./utils/docker/removeExistingNetwork.js";
import {readLatestConfig} from "./utils/config/readLatestConfig.js";

async function main() {

    const defaultCommand = program
        .name('xnetwork')
        .version('0.0.1')
        .description('An all-in-one tool for creating and managing your own MultiversX network')
        .option('--no-cache', 'Do not use cache when downloading files')
        .option('--custom-repo-path <path>', 'Fetch files from a custom repository')
        .option('--custom-repo-branch <branch>', 'Fetch files from a custom repository branch')
        .action(startInteractiveSetup)

    defaultCommand
        .command('create')
        .description('Create a new network, non-interactive')
        .argument('<config-path>', 'Path to the config file')
        .action((configPath) => {createNetworkAction(configPath)})

    defaultCommand
        .command('remove')
        .description('Remove a network, non-interactive')
        .action(removeNetworkAction)

    const configCommand = defaultCommand
        .command('config')
        .description('Config utils')

    configCommand
        .command('generate')
        .description('Generate a config file')
        .argument('<output-path>', 'Path to the output file')
        .action((outputPath) => {generateConfigAction(outputPath)})

    configCommand
        .command('latest')
        .description('Get the latest config file used to create a network, if exists')
        .action(getLatestConfigAction)

    await program.parse(process.argv)
}

export async function startInteractiveSetup() {
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

export async function getLatestConfigAction() {
    try {
        const latestConfig = await readLatestConfig()
        console.log(JSON.stringify(latestConfig, null, 4))
    } catch (e) {
        console.log('No latest config found')
    }
}

main().then(() => console.log('Done'))
