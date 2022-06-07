// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract SimpleCoin {
    mapping (address => uint256) public coinBalance;
    mapping (address => mapping (address => uint256)) public allowance;

    // =========================================================================

    event Log(string value);
    event Transfer(address indexed from, address indexed to, uint256 value);

    // =========================================================================

    constructor(uint256 initialSupply) {
        coinBalance[msg.sender] = initialSupply;
    }

    // =========================================================================

    function transfer(address to, uint256 amount) external {

        // Check that the sender has enough coins.
        if (coinBalance[msg.sender] < amount) {
            string memory resp = string(abi.encodePacked("insufficent funds  bal:", uint2str(coinBalance[msg.sender]), "  amount: ", uint2str(amount)));
            revert(resp);
        }

        // Check the amount is not zero or the amount value caused the unsigned
        // int to restart back to zero.
        if (coinBalance[to]+amount <= coinBalance[to]) {
            string memory resp = string(abi.encodePacked("invalid amount: ", uint2str(amount)));
            revert(resp);
        }

        emit Log("starting transfer");

        coinBalance[msg.sender] -= amount;
        coinBalance[to] += amount;

        emit Log("ending transfer");

        emit Transfer(msg.sender, to, amount);
    }

    // authorize allows the sender of this call to be authorized to spend the
    // specified accounts money up to the allowance amount.
    function authorize(address authorizedAccount, uint256 allowanceAmount) public returns (bool) {
        allowance[msg.sender][authorizedAccount] = allowanceAmount;
        return true;
    }

    // =========================================================================

    function uint2str(uint i) internal pure returns (string memory) {
        if (i == 0) {
            return "0";
        }

        uint j = i;
        uint len;
        while (j != 0) {
            len++;
            j /= 10;
        }

        bytes memory bstr = new bytes(len);
        uint k = len;
        while (i != 0) {
            k = k-1;
            uint8 temp = (48 + uint8(i - i / 10 * 10));
            bytes1 b1 = bytes1(temp);
            bstr[k] = b1;
            i /= 10;
        }

        return string(bstr);
    }
}