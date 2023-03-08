import * as os from "os"
import * as path from "path"
import { program } from "commander"

export class Constants {
    static CLI_USER_PATH = path.join(os.homedir(), '.xnetwork')
    static get REPO_URL(): string {
        return program.opts().customRepoPath ?? "git@github.com:gfusee/xnetwork.git"
    }

    static get REPO_BRANCH(): string {
        return program.opts().customRepoBranch ?? "main"
    }

    static get CLI_USER_REPO_PATH(): string {
        return path.join(Constants.CLI_USER_PATH, '.repo')
    }
}

