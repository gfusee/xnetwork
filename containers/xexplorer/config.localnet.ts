import { NetworkType } from 'types/network.types';
import { allApps, schema } from './sharedConfig';
export * from './sharedConfig';

export const networks: NetworkType[] = [
    {
        default: true,
        id: 'localnet',
        name: 'Localnet',
        chainId: 'localnet',
        adapter: 'api',
        theme: 'testnet',
        egldLabel: 'tEGLD',
        walletAddress: 'https://devnet-wallet.multiversx.com',
        explorerAddress: 'http://localhost:3002',
        nftExplorerAddress: 'https://devnet.xspotlight.com',
        apiAddress: 'http://localhost:3001'
    }
];

export const multiversxApps = allApps([
    {
        id: 'wallet',
        url: 'https://devnet-wallet.multiversx.com'
    },
    {
        id: 'explorer',
        url: 'http://localhost:3002'
    },
    {
        id: 'xexchange',
        url: 'http://devnet.xexchange.com'
    },
    {
        id: 'xspotlight',
        url: 'https://devnet.xspotlight.com/'
    }
]);

networks.forEach((network) => {
    schema.validate(network, { strict: true }).catch(({ errors }) => {
        console.error(`Config invalid format for ${network.id}`, errors);
    });
});
