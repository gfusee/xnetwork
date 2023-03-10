import {ExecException, ExecOptions} from "child_process"
import { exec } from "child_process"
import { createInterface } from 'readline'
import {Constants} from "../config/constants.js";

export type ExecError = {
    error: ExecException,
    stdout: string | Buffer
    stderr: string | Buffer
}

export function execCustom(
    command: string,
    logStdout?: boolean,
    options?: ExecOptions
): Promise<{stdout: string | Buffer, stderr: string | Buffer}> {
    return new Promise((resolve, reject) => {
        const child = exec(command, options,(error, stdout, stderr) => {
            if (error) {
                const thrownError: ExecError = {
                    error: error,
                    stdout: stdout,
                    stderr: stderr
                }
                reject(thrownError)
            } else {
                resolve({
                    stdout: stdout,
                    stderr: stderr
                })
            }
        })

        if (logStdout) {
            const reader = createInterface(child.stdout!)
            reader.on('line', (line) => {
                console.log(line)
            })
        }
    })

}

export function execCustomInRepo(
    command: string,
    logStdout?: boolean,
    options?: ExecOptions
): Promise<{stdout: string | Buffer, stderr: string | Buffer}> {
    return execCustom(command, logStdout, {
        cwd: Constants.CLI_USER_REPO_PATH,
        ...options
    })
}
