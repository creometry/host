// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import '@chainlink/contracts/src/v0.8/ChainlinkClient.sol';
import '@chainlink/contracts/src/v0.8/ConfirmedOwner.sol';

contract TokenFarm is ChainlinkClient, ConfirmedOwner {
    using Chainlink for Chainlink.Request;
    event RequestCoinsCalculation(bytes32 indexed requestId, uint256 coins);
    string public name = "Creometry community cluster";
    IERC20 public creoToken;
    bytes32 private jobId;
    uint256 private fee;
    uint256 public volume;
    event RequestVolume(bytes32 indexed requestId, uint256 volume);

    address[] public providers;
    mapping(address => bool) public hasProvided;
    mapping(address => bool) public isProviding;
    mapping(bytes32 => address) requesters;

    constructor(IERC20 _creoToken)  ConfirmedOwner(msg.sender) {
        creoToken = IERC20(_creoToken);
        setChainlinkToken(0x01BE23585060835E02B77ef475b0Cc51aA1e0709);
        setChainlinkOracle(0xf3FBB7f3391F62C8fe53f89B41dFC8159EE9653f);
        jobId = 'ca98366cc7314957b8c012c72f05aeeb';
        fee = (1 * LINK_DIVISIBILITY) / 10; // 0,1 * 10**18 (Varies by network and job)
    }

    function ProvideRessource() public {

        // before this we will check our api to test if this address really exists or no.
        if(!hasProvided[msg.sender]) {
            providers.push(msg.sender);
        }

        // Update Providing status
        isProviding[msg.sender] = true;
        hasProvided[msg.sender] = true;
    }

    function StopProviding() public {

        // Update Provider status
        isProviding[msg.sender] = false;
    }
    function toAsciiString(address x) internal pure returns (string memory) {
        bytes memory s = new bytes(40);
        for (uint i = 0; i < 20; i++) {
            bytes1 b = bytes1(uint8(uint(uint160(x)) / (2**(8*(19 - i)))));
            bytes1 hi = bytes1(uint8(b) / 16);
            bytes1 lo = bytes1(uint8(b) - 16 * uint8(hi));
            s[2*i] = char(hi);
            s[2*i+1] = char(lo);            
        }
        return string(s);
    }
    function char(bytes1 b) internal pure returns (bytes1 c) {
        if (uint8(b) < 10) return bytes1(uint8(b) + 0x30);
        else return bytes1(uint8(b) + 0x57);
}
    function concatenate(string memory s1, string memory s2) public pure returns (string memory) {
        return string(abi.encodePacked(s1, s2));
    }

    function harvest() public returns (bytes32 requestId) {
        //1) check if it's a valid address (msg.sender in isProviding.)

        //2)launch the request.
        Chainlink.Request memory req = buildChainlinkRequest(jobId, address(this), this.fulfill.selector);
        string memory url="https://afternoon-headland-32456.herokuapp.com/wallets/";

        // Set the URL to perform the GET request on
        req.add('get', concatenate(url,toAsciiString(msg.sender)));

        // Set the path to find the desired data in the API response, where the response format is:
        // {"RAW":
        //   {"ETH":
        //    {"USD":
        //     {
        //      "VOLUME24HOUR": xxx.xxx,
        //     }
        //    }
        //   }
        //  }
        // request.add("path", "RAW.ETH.USD.VOLUME24HOUR"); // Chainlink nodes prior to 1.0.0 support this format
        req.add('path', "coins"); // Chainlink nodes 1.0.0 and later support this format

        // Multiply the result by 1000000000000000000 to remove decimals
        int256 timesAmount = 10**18;
        req.addInt('times', timesAmount);

        // Sends the request
        bytes32 _requestId = sendChainlinkRequest(req,fee);
        requesters[_requestId]=msg.sender;
        return _requestId;
    }
    function fulfill(bytes32 _requestId, uint256 _volume) public recordChainlinkFulfillment(_requestId) {
        emit RequestVolume(_requestId, _volume);
        // transferring the amount that came from the server.
        creoToken.transfer(requesters[_requestId], _volume);
        
    }
    /**
    * Allow withdraw of Link tokens from the contract
     */
    function withdrawLink() public onlyOwner {
        LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
        require(link.transfer(msg.sender, link.balanceOf(address(this))), 'Unable to transfer');
    }
    function getVolume() public view returns(uint256) {
        return volume;
    }

}