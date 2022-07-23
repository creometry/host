import React, { Component } from 'react';
import Web3 from 'web3';
import CreoToken from '../abis/CreoToken.json';
import TokenFarm from '../abis/TokenFarm.json';
import Navbar from './Navbar';
import Main from './Main';
import './App.css';

class App extends Component {
  async componentWillMount() {
    await this.loadWeb3();
    await this.loadBlockchainData();
    //await this.loadData(this.state.account);
  }

  loadDataFunc = async (account) => {
    let loadData = async (userAccount) => {
      try {
        const res = await fetch(
          'https://afternoon-headland-32456.herokuapp.com/wallets/' +
            userAccount
        );
        const jsonRes = await res.json();
        console.log(jsonRes.coins);
        this.setState({
          TheRewardFromLastHarverst: jsonRes.coins,
        });
      } catch (e) {
        console.log(e);
      }
    };
    await loadData(account);
  };

  async loadBlockchainData() {
    const web3 = window.web3;

    const accounts = await web3.eth.getAccounts();
    this.setState({ account: accounts[0] });

    const networkId = await web3.eth.net.getId();

    // Load CreoToken
    const creometryData = CreoToken.networks[networkId];
    if (creometryData) {
      const creoToken = new web3.eth.Contract(
        CreoToken.abi,
        creometryData.address
      );
      console.log(creometryData.address);
      this.setState({ creoToken });
      let CreoBalance = await creoToken.methods
        .balanceOf(this.state.account)
        .call();
      this.setState({ CreoBalance: CreoBalance.toString() });
    } else {
      window.alert('DappToken contract not deployed to detected network.');
    }

    // Load TokenFarm
    const tokenFarmData = TokenFarm.networks[networkId];
    if (tokenFarmData) {
      const tokenFarm = new web3.eth.Contract(
        TokenFarm.abi,
        tokenFarmData.address
      );
      this.setState({ tokenFarm });
    } else {
      window.alert('TokenFarm contract not deployed to detected network.');
    }

    this.setState({ loading: false });
  }

  async loadWeb3() {
    if (window.ethereum) {
      window.web3 = new Web3(window.ethereum);
      await window.ethereum.enable();
    } else if (window.web3) {
      window.web3 = new Web3(window.web3.currentProvider);
    } else {
      window.alert(
        'Non-Ethereum browser detected. You should consider trying MetaMask!'
      );
    }
  }
  harvest = () => {
    this.setState({ loading: true });
    this.state.tokenFarm.methods
      .harvest()
      .send({ from: this.state.account })
      .on('transactionHash', (hash) => {
        this.setState({ loading: false });
      });
  };
  constructor(props) {
    super(props);
    this.state = {
      account: '0x0',
      tokenFarm: {},
      CreoBalance: '0',
      TheRewardFromLastHarverst: '0',
      loading: true,
    };
  }

  render() {
    let content;
    if (this.state.loading) {
      content = (
        <p id="loader" className="text-center">
          Loading...
        </p>
      );
    } else {
      content = (
        <Main
          CreoBalance={this.state.CreoBalance}
          account={this.state.account}
          coins={this.state.TheRewardFromLastHarverst}
          loadData={this.loadDataFunc}
          harvest={this.harvest}
        />
      );
    }

    return (
      <div>
        <Navbar account={this.state.account} />
        <div className="container-fluid mt-5">
          <div className="row">
            <main
              role="main"
              className="col-lg-12 ml-auto mr-auto"
              style={{ maxWidth: '600px' }}
            >
              <div className="content mr-auto ml-auto">
                <a
                  href="https://creometry.com/en"
                  target="_blank"
                  rel="noopener noreferrer"
                ></a>

                {content}
              </div>
            </main>
          </div>
        </div>
      </div>
    );
  }
}

export default App;
