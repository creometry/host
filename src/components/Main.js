import React, { Component } from 'react';
import creoLogo from '../creoLogo.png';

class Main extends Component {
  render() {
    return (
      <div id="content" className="mt-3">
        <table className="table table-borderless text-muted text-center">
          <thead>
            <tr>
              <th scope="col">
                Your reward
                <button
                  type="submit"
                  className="btn btn-link btn-block btn-sm"
                  onClick={(event) => {
                    event.preventDefault();
                    this.props.loadData(this.props.account);
                  }}
                >
                  Give me my current unharvested rewards
                </button>
              </th>
              <th scope="col">Your Current Creo Balance</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>
                here goes the live reward(still dummy data): {this.props.coins}
              </td>
              <td>
                {window.web3.utils.fromWei(this.props.CreoBalance, 'Ether')}{' '}
                CREO TOKEN
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    );
  }
}

export default Main;
