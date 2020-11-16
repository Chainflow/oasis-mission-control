![Image description](https://github.com/chris-remus/cosmos-validator-mission-control/blob/master/Untitled%20design.jpg)

# Oasis Mission Control Validator Monitoring and Alerting Dashboard 

## Background

[Chainflow](https://chainflow.io/staking) and our development partner [Vitwit](https://vitwit.com) have been awarded an Oasis grant to build the Oasis Mission Control Validator Monitoring and Alerting Dashboard. 

Together, we've [built a similar tool](https://chainflow.io/cosmos-validator-mission-control/) for the Cosmos community, under an Interchain Foundation Grant.

## Current Status - 16 November 2020 - Code Released for Use

You can find -

1 - Dashboard details, screenshots and alerts list [here](https://github.com/Chainflow/oasis-mission-control/issues/3)

2 - Code and installation instructions [here](./INSTRUCTIONS.md)

Have feedback or questions? Open an issue in this repo, thanks!

## Problem

The Oasis protocol requires validators to secure the network. These validators need to meet strict and demanding availability and security requirements. Monitoring and alerting tools are needed to ensure these validators are operating as expected. Using these tools helps validator operators provide the availability and security required and demanded by the network and its users.
 
Establishing these tools as a helpful and necessary infrastructure component before mainnet launch establishes good practices among validators. Doing so now makes it more likely these tools get implemented, before the additional demands and distractions of a live production network consume an increasing amount of validator time and attention.
 
Open source monitoring and alerting tools offer a partial solution. Yet these tools need customization to monitor and alert on key validator functions. Developing these tools often takes a back seat to "fire-fighting" activities. These activities ofen hijack a validator operator's daily operations. As such, these tools become seen as a "luxury" to "get to later", and as such don't get developed or implemented. Furthermore, smaller validator operators, key to stake decentralization, may not have the resources required to develop such tools.

## Solution

This project will use open source tools as its foundation. We will customize the tools, building the necessary plugins, etc., to provide a more comprehensive, Oasis-validator-specific monitoring and alerting toolset. We will provide the source code and documentation necessary for validators to implement the solution for themselves.
 
This will save them considerable time, attention and effort, making it much more likely they will establish this important piece of a highly available and secure validator operation. As a result, the Oasis network will benefit from a set of more reliable and secure validators, who are able to prevent and address infrastructure issues, before they become major issues.
 
This should free operator time and attention to put toward other value-added ecosystem activities, e.g. thoughtfully participating in network governance.

## For the Oasis Community

We're building this for the Oasis Community. While the Oasis validators may seem like the most obvious beneficiaries, the broader Oasis community will also benefit from a more secure network, more diverse set of validators and greater validator governance participation.

There will be a number of checkpoints along the way to collect community feedback. We look forward to the community's feedback, to help align what's built to the specific needs of the Oasis community.
