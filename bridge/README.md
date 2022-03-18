# cBridge between BAS and BSC
To bridge assets other than BAS native gas token, operators should setup and run [cBridge](https://cbridge.celer.network/) solution between BAS and BSC.
The solution includes onchain contracts on BAS and BSC, and a dedicated [State Guardian Network](https://cbridge-docs.celer.network/introduction/state-guardian-network) to ensure security. Note the validators for SGN can be different from BAS's validator set.

## Initial onchain setup
1. On BAS, deploy BEP20 token contract or use an existing one (eg. wrap of native gas token) as cBridge staking token, referred to as *cBridgeStakeToken*
2. On BAS, deploy Celer SGN core contracts
3. On both BAS and BSC, deploy separately cBridge contracts

See details at [1_onchain.md](./1_onchain.md)

## Start dedicated SGN


## cBridge operation
### Add new token
### Add new SGN validator
