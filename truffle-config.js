require('babel-register');
require('babel-polyfill');
const HDWalletProvider = require('@truffle/hdwallet-provider');
const mnemonic = process.env.MNEMONIC;
const url = process.env.RPC_URL;
module.exports = {
  networks: {
    development: {
      host: '127.0.0.1',
      port: 7545,
      network_id: '*', // Match any network id
    },
    rinkeby: {
      provider: () => {
        return new HDWalletProvider(mnemonic, url);
      },
      network_id: 4,
      skipDryRun: true,
    },
  },
  contracts_directory: './src/contracts/',
  contracts_build_directory: './src/abis/',
  compilers: {
    solc: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
      evmVersion: 'petersburg',
      version: '0.8.7',
    },
  },
};
