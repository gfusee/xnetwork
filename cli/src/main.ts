#!/usr/bin/env node

import { program } from 'commander'
import {
    createNetworkAction,
    generateConfigAction,
    getLatestConfigAction,
    removeNetworkAction,
    startInteractiveAction
} from "./actions";

async function main() {

    const defaultCommand = program
        .name('xnetwork')
        .version('0.0.1')
        .description('An all-in-one tool for creating and managing your own MultiversX network')
        .option('--no-cache', 'Do not use cache when downloading files')
        .option('--custom-repo-path <path>', 'Fetch files from a custom repository')
        .option('--custom-repo-branch <branch>', 'Fetch files from a custom repository branch')
        .action(startInteractiveAction)

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

main().then(() => console.log('Done'))
