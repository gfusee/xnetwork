export type CLIConfig = {
    shouldHaveElasticSearch: boolean
    shouldHaveMySQL: boolean
    shouldHaveRedis: boolean
    shouldHaveRabbitMQ: boolean
    shouldHaveApi: boolean
    shouldHaveXExplorer: boolean
    numberOfShards: number
    initialEGLDAddress?: string
    mxOpsScenesPath?: string
}

export function getDefaultConfig(): CLIConfig {
    return {
        shouldHaveElasticSearch: false,
        shouldHaveMySQL: false,
        shouldHaveRedis: false,
        shouldHaveRabbitMQ: false,
        shouldHaveApi: false,
        shouldHaveXExplorer: false,
        numberOfShards: 1
    }
}
