import * as os from "os"
import * as path from "path";

export class Constants {
    static CLI_USER_PATH = path.join(os.homedir(), '.xnetwork')
    static REPO_URL = "~/IdeaProjects/mx-full-testnet-docker"
    static get CLI_USER_REPO_PATH(): string {
        return path.join(Constants.CLI_USER_PATH, '.repo')
    }
}

