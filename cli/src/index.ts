#!/usr/bin/env node

import {getDefaultConfig} from "./config/config.js"
import {RunnerQuestion} from "./questions/runner/runnerQuestion.js"
import {ResultLogger} from "./result/resultLogger.js"
import {NumberShardsQuestion} from "./questions/fresh/numberShardsQuestion.js"
import {DockerPrerequisites} from "./prerequisites/dockerPrerequisites.js"
import {dontIndent} from "./utils/strings/dontIndent.js";
import chalk from "chalk";

async function main() {
    console.log('Checking prerequisites...')

    try {
        await (new DockerPrerequisites()).check()
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
    await (new NumberShardsQuestion()).process(config)
    await (new RunnerQuestion()).process(config)

    const resultLogger = new ResultLogger()
    await resultLogger.printResults(config)
}

main().then(() => console.log('Done'))
