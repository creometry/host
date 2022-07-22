const CreoToken = artifacts.require('CreoToken');
const TokenFarm = artifacts.require('TokenFarm');

module.exports = async function(deployer, network, accounts) {
  // Deploy Dapp Token
  await deployer.deploy(CreoToken);
  const creoToken = await CreoToken.deployed();

  // Deploy TokenFarm
  await deployer.deploy(TokenFarm, creoToken.address);
  const tokenFarm = await TokenFarm.deployed();

  // Transfer all tokens to TokenFarm (1 million)
  //await dappToken.transfer(tokenFarm.address, '1000000000000000000000000');

  // Transfer 100 Mock DAI tokens to investor
  //await daiToken.transfer(accounts[1], '100000000000000000000')
};
