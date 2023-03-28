import * as os from "os"
import * as path from "path"
import { program } from "commander"
import fs from "fs"

export enum PauseBehavior {
    PAUSE,
    STOP
}

export type ContainerInfos = {
    name: string
    pauseBehavior: PauseBehavior
}

export class Constants {
    static CLI_USER_PATH = path.join(os.homedir(), '.xnetwork')
    static get REPO_URL(): string {
        return program.opts().customRepoPath ?? "https://github.com/gfusee/xnetwork.git"
    }

    static get REPO_BRANCH(): string | undefined {
        return program.opts().customRepoBranch ?? `v${Constants.getPackageJson().version}`
    }

    static get CLI_USER_REPO_PATH(): string {
        return path.join(Constants.CLI_USER_PATH, '.repo')
    }

    static get CLI_USER_STORAGE_PATH(): string {
        return path.join(Constants.CLI_USER_PATH, 'storage')
    }

    static get CLI_USER_STORAGE_LATEST_CONFIG(): string {
        return path.join(Constants.CLI_USER_STORAGE_PATH, 'latest_config.json')
    }

    static get TESTNET_CONTAINER(): ContainerInfos {
        return {
            name: "testnet",
            pauseBehavior: PauseBehavior.PAUSE
        }
    }

    static get API_CONTAINER(): ContainerInfos {
        return {
            name: "api",
            pauseBehavior: PauseBehavior.STOP
        }
    }

    static get ELASTIC_CONTAINER(): ContainerInfos {
        return {
            name: "elastic",
            pauseBehavior: PauseBehavior.PAUSE
        }
    }

    static get MYSQL_CONTAINER(): ContainerInfos {
        return {
            name: "mysql",
            pauseBehavior: PauseBehavior.STOP
        }
    }

    static get REDIS_CONTAINER(): ContainerInfos {
        return {
            name: "redis",
            pauseBehavior: PauseBehavior.STOP
        }
    }

    static get RABBITMQ_CONTAINER(): ContainerInfos {
        return {
            name: "rabbitmq",
            pauseBehavior: PauseBehavior.STOP
        }
    }

    static isNoCacheRequested(): boolean {
        return !(program.opts().cache ?? true)
    }

    static getPackageJson(): { version: string } {
        return JSON.parse(fs.readFileSync(path.join(__dirname, '../../package.json'), 'utf-8'))
    }
}

