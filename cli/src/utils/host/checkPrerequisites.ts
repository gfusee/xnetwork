import {DockerPrerequisites} from "../../prerequisites/dockerPrerequisites.js";
import {GitRepoPrerequisites} from "../../prerequisites/gitRepoPrerequisites.js";

export async function checkPrerequisites() {
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
}
