#!/usr/bin/env node

import {getDefaultConfig} from "./config/config.js"
import {DockerPrerequisites} from "./prerequisites/dockerPrerequisites.js"
import {dontIndent} from "./utils/strings/dontIndent.js";
import chalk from "chalk";
import {GitRepoPrerequisites} from "./prerequisites/gitRepoPrerequisites.js"
import { program } from 'commander'
import {StartQuestion} from "./questions/startQuestion.js";

async function main() {

    program
        .name('xnetwork')
        .version('0.0.1')
        .description('An all-in-one tool for creating and managing your own MultiversX network')
        .option('--no-cache', 'Do not use cache when downloading files')
        .option('--custom-repo-path <path>', 'Fetch files from a custom repository')
        .option('--custom-repo-branch <branch>', 'Fetch files from a custom repository branch')

    await program.parse(process.argv)

    console.log('Checking prerequisites...')

    try {
        await (new DockerPrerequisites()).check()
        await (new GitRepoPrerequisites()).check()
    } catch (e) {
        if (typeof e === 'string') {
            console.log(e)
            process.exit(1)
        } else {
            throw e
        }
    }

    const welcomeMessage = `
    Welcome to ${chalk.bold('xNetwork')} CLI! 🔥
    
    This CLI will help you set up a network with the features you want.
    Let's start by asking you a few questions.
    `

    console.log(dontIndent(welcomeMessage))

    const config = getDefaultConfig()
    await (new StartQuestion()).process(config)
}

main().then(() => console.log('\nDone'))
