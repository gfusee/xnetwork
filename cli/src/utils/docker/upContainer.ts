import {execCustomInRepo} from "../exec.js";

export async function upContainer(containerName: string, env?: NodeJS.ProcessEnv) {
    await execCustomInRepo(`docker-compose build --no-cache ${containerName}`, false, {env: env})
    await execCustomInRepo(`docker-compose up -d ${containerName}`, false, {env: env})
}
