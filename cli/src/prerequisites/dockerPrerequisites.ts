import {Prerequisites} from "./prerequisites.js"
import {execCustom} from "../utils/exec.js"
import ora from "ora"
import chalk from 'chalk'

const dockerNotInstalledMessage = `
${chalk.red.bold('Sorry, it looks like Docker is not installed on your system.\n\n')}
To use this CLI tool, you'll need to install Docker. Please follow the instructions below for your operating system:

${chalk.blue.bold('- For macOS:')}
    1. Download the Docker Desktop installer from ${chalk.underline('https://www.docker.com/products/docker-desktop')}.
    2. Double-click the downloaded .dmg file to start the installation process.
    3. Follow the prompts to complete the installation.

${chalk.blue.bold('- For Windows:')}
    1. Download the Docker Desktop installer from ${chalk.underline('https://www.docker.com/products/docker-desktop')}.
    2. Double-click the downloaded .exe file to start the installation process.
    3. Follow the prompts to complete the installation.

${chalk.blue.bold('- For Linux:')}
    1. Please follow the instructions specific to your Linux distribution on the Docker website: ${chalk.underline('https://docs.docker.com/engine/install/')}

Once Docker is installed, please run this CLI tool again. Thank you!
`

const dockerDaemonNotRunningMessage = `
${chalk.red.bold('Sorry, it looks like the Docker Daemon is not running on your system.\n\n')}
To use this CLI tool, you'll need to start Docker first. Please follow the instructions below for your operating system:

${chalk.blue.bold('- For macOS:')}
    1. Open the "Docker Desktop" application from your Applications folder.
    2. Wait for Docker to start. You can check its status in the menu bar.
    3. Once Docker is running, please run this CLI tool again.

${chalk.blue.bold('- For Windows:')}
    1. Open the "Docker Desktop" application from your Start menu.
    2. Wait for Docker to start. You can check its status in the system tray.
    3. Once Docker is running, please run this CLI tool again.

${chalk.blue.bold('- For Linux:')}
    1. Please follow the instructions specific to your Linux distribution on the Docker website: ${chalk.underline('https://docs.docker.com/engine/install/')}.
    2. Once Docker is installed, start it by running the command "sudo systemctl start docker" in your terminal.
    3. Once Docker is running, please run this CLI tool again.

If you have any questions or need further assistance, please consult the Docker documentation or community resources. Thank you!
`

export class DockerPrerequisites implements Prerequisites {

    async check() {
        const dockerSpinner = ora('Checking if docker is installed...').start()
        try {
            await execCustom('which docker')
        } catch (e) {
            dockerSpinner.fail('Docker is not installed')
            throw dockerNotInstalledMessage
        }
        dockerSpinner.succeed('Docker is installed')

        const dockerComposeSpinner = ora('Checking if docker compose is installed...').start()
        try {
            await execCustom('docker compose')
        } catch (e) {
            dockerComposeSpinner.fail('Docker compose is not installed')
            throw dockerNotInstalledMessage
        }
        dockerComposeSpinner.succeed('Docker compose is installed')

        const dockerDaemonSpinner = ora('Checking if docker daemon is running...').start()
        try {
            await execCustom('docker info')
        } catch (e) {
            dockerDaemonSpinner.fail('Docker daemon is not running')
            throw dockerDaemonNotRunningMessage
        }
        dockerDaemonSpinner.succeed('Docker daemon is running')
    }

}
