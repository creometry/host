// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract CreoToken is ERC20 {
    constructor() ERC20("CREO Token", "CRT") {
        _mint(msg.sender, 1000000*10**18);
    }
}
