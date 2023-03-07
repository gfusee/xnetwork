export type CLIConfig = {
    shouldHaveElasticSearch: boolean
    shouldHaveMySQL: boolean
    shouldHaveRedis: boolean
    shouldHaveRabbitMQ: boolean
    shouldHaveApi: boolean
    numberOfShards: number
    initialEGLDAddress?: string
}

export function getDefaultConfig(): CLIConfig {
    return {
        shouldHaveElasticSearch: false,
        shouldHaveMySQL: false,
        shouldHaveRedis: false,
        shouldHaveRabbitMQ: false,
        shouldHaveApi: false,
        numberOfShards: 1
    }
}
