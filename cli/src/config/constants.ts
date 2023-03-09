import * as os from "os"
import * as path from "path"
import { program } from "commander"
import fs from "fs"

import { fileURLToPath } from 'url'
import { dirname } from 'path'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

export class Constants {
    static CLI_USER_PATH = path.join(os.homedir(), '.xnetwork')
    static get REPO_URL(): string {
        return program.opts().customRepoPath ?? "git@github.com:gfusee/xnetwork.git"
    }

    static get REPO_BRANCH(): string | undefined {
        return program.opts().customRepoBranch ?? `v${Constants.getPackageJson().version}`
    }

    static get CLI_USER_REPO_PATH(): string {
        return path.join(Constants.CLI_USER_PATH, '.repo')
    }

    static isNoCacheRequested(): string {
        return program.opts().cache ?? true
    }

    static getPackageJson(): { version: string } {
        return JSON.parse(fs.readFileSync(path.join(__dirname, '../../package.json'), 'utf-8'))
    }
}

