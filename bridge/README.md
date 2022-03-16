# cBridge between BAS and BSC
Setup and run [cBridge](https://cbridge.celer.network/) solution between BAS and BSC, as BSC's native asset bridge only supports gas token.
The solution includes onchain contracts on BAS and BSC, and a dedicated [State Guardian Network](https://cbridge-docs.celer.network/introduction/state-guardian-network) to ensure security.
Note the validators for SGN can be different from BSC's validator set.

## Initial onchain setup
1. On BAS, deploy ERC20 token contract or use an existing one (eg. ERC20 wrap of gas token) as cBridge staking token, referred to as *StakeToken*
2. On BAS, deploy Celer SGN core contracts
3. On BAS and BSC, deploy cBridge contracts and setup signers

## Start dedicated SGN


## cBridge operation
### Add new token
### Add new SGN validator
