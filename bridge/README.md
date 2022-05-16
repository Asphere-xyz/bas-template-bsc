# cBridge between BAS and BSC

To bridge assets other than the BAS native gas token, operators can choose to bridge BEP-20 tokens via
[cBridge](https://cbridge.celer.network/). It is also possible to setup and run a standalone cBridge solution between BAS and BSC if the
operators choose to do so. Note that this is a technically advanced topic that requires good understanding of both Ethereum based and Cosmos
SDK based blockchains. The guide is still a work in progress and might contain errors or steps that are unclear.

The standalone solution includes on-chain contracts on BAS and BSC, and a dedicated [State Guardian Network](https://cbridge-docs.celer.network/introduction/state-guardian-network) as the relay chain. Note the validators for SGN can be different from those of the BAS.

## Steps

Once you have followed the steps below and got the SGN chain up and running, **contact the Celer team** to get a dedicated gateway server and
frontend setup for your BAS.

- [Deploy Contracts](01_deploy_contracts.md)
- [Init SGN](02_init_sgn.md)
- [Run Validator](03_run_validator.md)
- [Add Chains and Tokens](04_add_chain_tokens.md)
