import axios from "axios"

export async function waitForAPIToBeReady() {
    const url = "http://localhost:3001/hello"

    while (true) {
        try {
            const response1 = await axios.get(url)
            if (response1.status === 200) {
                break
            }
        } catch (error) {
            // do nothing
        }

        await new Promise(resolve => setTimeout(resolve, 1000))
    }
}
