# METAYA

## Introduction

METAYA platform is aiming to create a new ecology of the digital content industry.

Recently, NFT has entered an explosive period. This greatly helps the entire digital content product market, and it has great power to lead the market to a new stage.

However, piracy and misappropriation in the digital content product market have become a prominent problem. So far, NFT, as the “rights certificate” on the blockchain is separated from the digital content product it represents, these two has no connection with each other.

METAYA utilizes SOUNDLINKS technology to implant indelible SOUNDLINKS DID into audio and video works, and registers SOUNDLINKS DID with the DRM information of audio and video products in NFT, so that the audio and video products offchain and the rights certificates onchain are anchored. When audio and video products are played back on traditional social media and content platforms, the SOUNDLINKS DID can be detected and directed to the blockchain to verify the rights. 

Many unique brands, contents and IPs have been joining with us in such a great march. We welcome more participation from all of you.

## What is Soundlinks DID ?
Verification protocol，Enabling NFT Protection.

Provides unique binding between NFT and digital asset associated:
1. Embeds "DID" into digital asset
1. "DID" Data simultaneously stored/anchored into the NFT minted
1. Uses and utilizes audio to transmit arbitrary info NFC tech that leverages audio and based on sound triggers

A same audio file can contain & transmit different information.

Soundlinks DID requires minimal storage ( no need to store audio content itself ) - significantly reduces on-chain storage, bandwidth & computing power for interaction between on-chain/off-chain.

- `Anti-piracy`: Off-chain digital content, played back anywhere off-chain is protected and verified, as the Soundlinks DID embedded in the digital content can be linked to the anchored NFT.

- `Digital content circulation`: Digital content can circulate freely off-chain ( existing infrastructure, such as content platforms, social platforms, etc. ). Only Soundlinks DID needs to be stored on-chain, not the digital content itself.

- `Low-carbon environmental protection, climate awareness`: Each Soundlinks DID <= 256bits; A MP3 format song is usually 5-6 Mbytes = 40-50 Mbits; Storage, bandwidth & computing power required by Soundlinks powered blockchain are about one in hundreds of thousands of storing files on the blockchain.

# METAYA Contract Addresses

`Metaya.cdc` : This is the main Metaya smart contract that defines the core functionality of the NFT, base on [NBA TopShot smart contract](https://github.com/dapperlabs/nba-smart-contracts).

| Network | Contract Address     |
|---------|----------------------|
| Testnet | [0x2c5bdf0e0d907421](https://flow-view-source.com/testnet/account/0x2c5bdf0e0d907421/contract/Metaya) |
| Mainnet | [0x8b935cd43003d4b2](https://flow-view-source.com/mainnet/account/0x8b935cd43003d4b2/contract/Metaya) |

`MetayaShardedCollection.cdc` : This contract bundles together a bunch of MomentCollection objects in a dictionary, and then distributes the individual Moments between them while implementing the same public interface as the default MomentCollection implementation.

| Network | Contract Address     |
|---------|----------------------|
| Testnet | [0x2c5bdf0e0d907421](https://flow-view-source.com/testnet/account/0x2c5bdf0e0d907421/contract/MetayaShardedCollection) |
| Mainnet | [0x8b935cd43003d4b2](https://flow-view-source.com/mainnet/account/0x8b935cd43003d4b2/contract/MetayaShardedCollection) |

`MetayaUtilityCoin.cdc` : The utility coins circulates on Metaya.

| Network | Contract Address     |
|---------|----------------------|
| Testnet | [0x2c5bdf0e0d907421](https://flow-view-source.com/testnet/account/0x2c5bdf0e0d907421/contract/MetayaUtilityCoin) |
| Mainnet | [0x8b935cd43003d4b2](https://flow-view-source.com/mainnet/account/0x8b935cd43003d4b2/contract/MetayaUtilityCoin) |

`MetayaBeneficiaryCut.cdc` : This smart contract stores the mappings from the names of copyright owners to the vaults in which they'd like to receive tokens, as well as the cut they'd like to take from store and pack sales revenue and marketplace transactions.

| Network | Contract Address     |
|---------|----------------------|
| Testnet | [0x2c5bdf0e0d907421](https://flow-view-source.com/testnet/account/0x2c5bdf0e0d907421/contract/MetayaBeneficiaryCut) |
| Mainnet | [0x8b935cd43003d4b2](https://flow-view-source.com/mainnet/account/0x8b935cd43003d4b2/contract/MetayaBeneficiaryCut)|

`NFTStorefront.cdc`: The general-purpose contract is used in the Metaya market.

| Network | Contract Address     |
|---------|----------------------|
| Testnet | [0x94b06cfca1d8a476](https://flow-view-source.com/testnet/account/0x94b06cfca1d8a476/contract/NFTStorefront) |
| Mainnet | [0x4eb8a10cb9f87357](https://flow-view-source.com/mainnet/account/0x4eb8a10cb9f87357/contract/NFTStorefront) |

# SOUNDLINKS Contract Address

`SoundlinksDID.cdc` : Each Metaya NFT is embedded with a unique Soundlinks DID.

| Network | Contract Address     |
|---------|----------------------|
| Testnet | [0x2771ed97c1150a08](https://flow-view-source.com/testnet/account/0x2771ed97c1150a08/contract/SoundlinksDID) |
| Mainnet | [0x602e888f32abc278](https://flow-view-source.com/mainnet/account/0x602e888f32abc278/contract/SoundlinksDID) |

# Common Commands

#### Deploy contract
```
flow project deploy --network=testnet
```
#### Update deployed contract
```
flow project deploy --network=testnet --update
```
#### Remove deployed contract
```
flow accounts remove-contract SoundlinksDID --network=testnet --signer=testnet-account-soundlinks

flow accounts remove-contract Metaya --network=testnet --signer=testnet-account-metaya
flow accounts remove-contract MetayaShardedCollection --network=testnet --signer=testnet-account-metaya
flow accounts remove-contract MetayaUtilityCoin --network=testnet --signer=testnet-account-metaya
flow accounts remove-contract MetayaBeneficiaryCut --network=testnet --signer=testnet-account-metaya
```

# Soundlinks DID Commands
SoundlinksDID contract is already deployed to testnet at [0x2771ed97c1150a08](https://flow-view-source.com/testnet/account/0x2771ed97c1150a08).

#### Setup Account `Transaction`
```
flow transactions send ./transactions/soundlinksDID/setup_account.cdc --signer testnet-account-metaya --network=testnet
```
#### Mint DIDs `Transaction`
```
flow transactions send ./transactions/soundlinksDID/mint_DIDs.cdc --signer testnet-account-soundlinks --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "Array","value": [{"type": "String","value": "b822bb93905a9bd8b3a0c08168c427696436cf8bf37ed4ab8ebf41a07642e111"},{"type": "String","value": "de689e8d537fd816753da1d4fa6873e16f8dfbbcfd5d9e5c9c35a0a426645222"}]}]'
```
#### Purchase DIDs `Transaction`
```
flow transactions build ./transactions/soundlinksDID/purchase_DIDs.cdc --network=testnet --args-json '[{"type": "UInt32","value": "2"},{"type": "Array","value": [{"type": "String","value": "b822bb93905a9bd8b3a0c08168c427696436cf8bf37ed4ab8ebf41a07642e555"},{"type": "String","value": "de689e8d537fd816753da1d4fa6873e16f8dfbbcfd5d9e5c9c35a0a426645666"}]},{"type": "UFix64","value": "1.0"}]' --authorizer testnet-account-soundlinks --authorizer testnet-account-metaya --proposer testnet-account-soundlinks --payer testnet-account-metaya --filter payload --save built.rlp

flow transactions sign ./built.rlp --signer testnet-account-soundlinks --network=testnet --filter payload --save signed.rlp

flow transactions sign ./signed.rlp --signer testnet-account-metaya --network=testnet --filter payload --save signed.rlp

flow transactions send-signed ./signed.rlp --network=testnet
```
#### Get supply `Script`
```
flow scripts execute ./scripts/soundlinksDID/get_supply.cdc --network=testnet
```
#### Get Amount `Script`
```
flow scripts execute ./scripts/soundlinksDID/get_amount.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"}]'
```

# Metaya Commands
Metaya contracts are already deployed to testnet at [0x2c5bdf0e0d907421](https://flow-view-source.com/testnet/account/0x2c5bdf0e0d907421).

## Metaya NFT

### Admin `Transaction`
---
#### Admin / Create Play `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/create_play.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "Dictionary","value": [{"key": {"type": "String","value": "Title"},"value": {"type": "String","value": "Play 001"}}]}]'
```
#### Admin / Create Set `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/create_set.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "String","value": "Set 001"}]'
```
#### Admin / Add Play to Set `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/add_play_to_set.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "UInt32","value": "1"}]'
```
#### Admin / Add Plays to Set `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/add_plays_to_set.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "4"},{"type": "Array","value": [{"type": "UInt32","value": "7"},{"type": "UInt32","value": "8"},{"type": "UInt32","value": "9"}]}]'
```
#### Admin / Start New Series `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/start_new_series.cdc --signer testnet-account-metaya --network=testnet
```
#### Admin / Lock Set `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/lock_set.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "2"}]'
```
#### Admin / Retire Play from Set `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/retire_play_from_set.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "3"},{"type": "UInt32","value": "4"}]'
```
#### Admin / Retire All Plays from Set `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/retire_allPlays_from_set.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "3"}]'
```
#### Admin / Mint Moment `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/mint_moment.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "UInt32","value": "1"},{"type": "Address","value": "0x2c5bdf0e0d907421"}]'
```
#### Admin / Batch Mint Moments `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/batch_mint_moments.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "4"},{"type": "UInt32","value": "7"},{"type": "UInt32","value": "2"},{"type": "Address","value": "0x2c5bdf0e0d907421"}]'
```
#### Admin / Provide Moment `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/provide_moment.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "Address","value": "0xd542317949eb00b6"},{"type": "Array","value": [{"type": "UInt64","value": "1"},{"type": "UInt64","value": "2"},{"type": "UInt64","value": "3"}]}]'
```
#### Admin / Purchase Store by Cash `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/purchase_store_by_cash.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "UFix64","value": "100.0"},{"type": "String","value": "Commonweal 001"}]'
```
#### Admin / Purchase Store by IUC `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/purchase_store_by_IUC.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "UFix64","value": "15.0"},{"type": "String","value": "Commonweal 001"}]'
```
#### Admin / Purchase Pack by Cash `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/purchase_pack_by_cash.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "UFix64","value": "100.0"},{"type": "String","value": "Commonweal 001"}]'
```
#### Admin / Purchase Pack by IUC `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/purchase_pack_by_IUC.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "UFix64","value": "15.0"},{"type": "String","value": "Commonweal 001"}]'
```
#### Admin / Transfer Admin `Transaction`
```
flow transactions send ./transactions/metayaNFT/admin/transfer_admin.cdc --signer testnet-account-metaya --network=testnet
```

### ShardedCollection `Transaction`
---
#### ShardedCollection / Setup Sharded Collection `Transaction`
```
flow transactions send ./transactions/metayaNFT/shardedCollection/setup_sharded_collection.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt64","value": "32"}]'
```
#### ShardedCollection / Transfer from Sharded `Transaction`
```
flow transactions send ./transactions/metayaNFT/shardedCollection/transfer_from_sharded.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "Address","value": "0xd542317949eb00b6"},{"type": "UInt64","value": "1"}]'
```
#### ShardedCollection / Batch Transfer from Sharded `Transaction`
```
flow transactions send ./transactions/metayaNFT/shardedCollection/batch_transfer_from_sharded.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "Address","value": "0xd542317949eb00b6"},{"type": "Array","value": [{"type": "UInt64","value": "2"},{"type": "UInt64","value": "3"}]}]'
```

### User `Transaction`
---
#### User / Setup Account `Transaction`
```
flow transactions send ./transactions/metayaNFT/user/setup_account.cdc --signer testnet-account --network=testnet
```
#### User / Transfer Moment `Transaction`
```
flow transactions send ./transactions/metayaNFT/user/transfer_moment.cdc --signer testnet-account --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"}]'
```
#### User / Batch Transfer Moments `Transaction`
```
flow transactions send ./transactions/metayaNFT/user/batch_transfer_moments.cdc --signer testnet-account --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "Array","value": [{"type": "UInt64","value": "2"},{"type": "UInt64","value": "3"}]}]'
```

### Collections `Script`
---
#### Collections / Get Collection IDs `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_collection_ids.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"}]'
```
#### Collections / Get ID in Collection `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_id_in_collection.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"}]'
```
#### Collections / Get Metadata `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_metadata.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"}]'
```
#### Collections / Get Metadata Field `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_metadata_field.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"},{"type": "String","value": "Title"}]'
```
#### Collections / Get Moment PlayID `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_moment_playID.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"}]'
```
#### Collections / Get Moment SerialNum `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_moment_serialNum.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"}]'
```
#### Collections / Get Moment Series `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_moment_series.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"}]'
```
#### Collections / Get Moment SetID `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_moment_setID.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"}]'
```
#### Collections / Get Moment SetName `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_moment_setName.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UInt64","value": "1"}]'
```
#### Collections / Get Set-Play are owned `Script`
```
flow scripts execute ./scripts/metayaNFT/collections/get_setplays_are_owned.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "Array","value": [{"type": "UInt32","value": "1"},{"type": "UInt32","value": "4"}]},{"type": "Array","value": [{"type": "UInt32","value": "3"},{"type": "UInt32","value": "9"}]}]'
```

### Plays `Script`
---
#### Plays / Get All Plays `Script`
```
flow scripts execute ./scripts/metayaNFT/plays/get_all_plays.cdc --network=testnet
```
#### Plays / Get Next PlayID `Script`
```
flow scripts execute ./scripts/metayaNFT/plays/get_nextPlayID.cdc --network=testnet
```
#### Plays / Get Play Metadata `Script`
```
flow scripts execute ./scripts/metayaNFT/plays/get_play_metadata.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"}]'
```
#### Plays / Get Play Metadata Field `Script`
```
flow scripts execute ./scripts/metayaNFT/plays/get_play_metadata_field.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "String","value": "Title"}]'
```

### Sets `Script`
---
#### Sets / Get Edition Retired `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_edition_retired.cdc --network=testnet --args-json '[{"type": "UInt32","value": "3"},{"type": "UInt32","value": "4"}]'
```
#### Sets / Get Next SetID `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_nextSetID.cdc --network=testnet
```
#### Sets / Get numMoments in edition `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_numMoments_in_edition.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "UInt32","value": "1"}]'
```
#### Sets / Get Plays in Set `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_plays_in_set.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"}]'
```
#### Sets / Get Set Series `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_setSeries.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"}]'
```
#### Sets / Get Set Name `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_setName.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"}]'
```
#### Sets / Get SetIDs by Name `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_setIDs_by_name.cdc --network=testnet --args-json '[{"type": "String","value": "Set 001"}]'
```
#### Sets / Get Set Locked `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_set_locked.cdc --network=testnet --args-json '[{"type": "UInt32","value": "2"}]'
```
#### Sets / Get Set Data `Script`
```
flow scripts execute ./scripts/metayaNFT/sets/get_set_data.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"}]'
```
#### Get Total Supply `Script`
```
flow scripts execute ./scripts/metayaNFT/get_totalSupply.cdc --network=testnet
```
#### Get Current Series `Script`
```
flow scripts execute ./scripts/metayaNFT/get_currentSeries.cdc --network=testnet
```

## Metaya Utility Coin

#### Setup Account `Transaction`
```
flow transactions send ./transactions/metayaUtilityCoin/setup_account.cdc --signer testnet-account-metaya --network=testnet
```
#### Mint Tokens `Transaction`
```
flow transactions send ./transactions/metayaUtilityCoin/mint_tokens.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"},{"type": "UFix64","value": "100.0"}]'
```
#### Transfer Tokens `Transaction`
```
flow transactions send ./transactions/metayaUtilityCoin/transfer_tokens.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UFix64","value": "50.0"},{"type": "Address","value": "0xd542317949eb00b6"}]'
```
#### Transfer Many Accounts `Transaction`
```
flow transactions send ./transactions/metayaUtilityCoin/transfer_many_accounts.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "Dictionary","value": [{"key": {"type": "Address","value": "0xd542317949eb00b6"},"value": {"type": "UFix64","value": "10.0"}},{"key": {"type": "Address","value": "0x312588a458110069"},"value": {"type": "UFix64","value": "10.0"}}]}]'
```
#### Burn Tokens by Admin `Transaction`
```
flow transactions send ./transactions/metayaUtilityCoin/burn_tokens_by_admin.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UFix64","value": "30.0"}]'
```
#### Burn Tokens by User `Transaction`
```
flow transactions build ./transactions/metayaUtilityCoin/burn_tokens_by_user.cdc --network=testnet --args-json '[{"type": "UFix64","value": "70.0"}]' --authorizer testnet-account --authorizer testnet-account-metaya --proposer testnet-account --payer testnet-account-metaya --filter payload --save built.rlp

flow transactions sign ./built.rlp --signer testnet-account --network=testnet --filter payload --save signed.rlp

flow transactions sign ./signed.rlp --signer testnet-account-metaya --network=testnet --filter payload --save signed.rlp

flow transactions send-signed ./signed.rlp --network=testnet
```
#### Get Supply `Script`
```
flow scripts execute ./scripts/metayaUtilityCoin/get_supply.cdc --network=testnet
```
#### Get Balance `Script`
```
flow scripts execute ./scripts/metayaUtilityCoin/get_balance.cdc --network=testnet --args-json '[{"type": "Address","value": "0x2c5bdf0e0d907421"}]'
```

## Metaya Beneficiary Cut

#### Set CopyrightOwner `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/set_copyrightOwner.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "String","value": "CopyrightOwner 001"},{"type": "Address","value": "0x312588a458110069"}]'
```
#### Del CopyrightOwner `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/del_copyrightOwner.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "String","value": "CopyrightOwner 002"}]'
```
#### Set Commonweal `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/set_commonweal.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "String","value": "Commonweal 001"},{"type": "Address","value": "0xc6834d0636ae584d"},{"type": "UFix64","value": "0.002"}]'
```
#### Del Commonweal `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/del_commonweal.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "String","value": "Commonweal 001"}]'
```
#### Set Metaya Capability `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/set_metaya_capability.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "Address","value": "0xd542317949eb00b6"}]'
```
#### Set Metaya Market CutPercentage `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/set_metaya_marketCutPercentage.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UFix64","value": "0.04"}]'
```
#### Set Store CutPercentage `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/set_storeCutPercentage.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "Dictionary","value": [{"key": {"type": "String","value": "Metaya"},"value": {"type": "UFix64","value": "0.2"}},{"key": {"type": "String","value": "CopyrightOwner 001"},"value": {"type": "UFix64","value": "0.3"}},{"key": {"type": "String","value": "CopyrightOwner 002"},"value": {"type": "UFix64","value": "0.5"}}]}]'
```
#### Del Store CutPercentage `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/del_storeCutPercentage.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"}]'
```
#### Set Pack CutPercentage `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/set_packCutPercentage.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "Dictionary","value": [{"key": {"type": "String","value": "Metaya"},"value": {"type": "UFix64","value": "0.3"}},{"key": {"type": "String","value": "CopyrightOwner 001"},"value": {"type": "UFix64","value": "0.1"}},{"key": {"type": "String","value": "CopyrightOwner 002"},"value": {"type": "UFix64","value": "0.6"}}]}]'
```
#### Del Pack CutPercentage `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/del_packCutPercentage.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"}]'
```
#### Set Market CutPercentage `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/set_marketCutPercentage.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "Dictionary","value": [{"key": {"type": "String","value": "CopyrightOwner 001"},"value": {"type": "UFix64","value": "0.03"}},{"key": {"type": "String","value": "CopyrightOwner 002"},"value": {"type": "UFix64","value": "0.07"}}]}]'
```
#### Del Market CutPercentage `Transaction`
```
flow transactions send ./transactions/metayaBeneficiaryCut/del_marketCutPercentage.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt32","value": "1"}]'
```
#### Get CopyrightOwner Names `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_copyrightOwner_names.cdc --network=testnet
```
#### Get CopyrightOwner Amount `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_copyrightOwner_amount.cdc --network=testnet
```
#### Get CopyrightOwner Contain `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_copyrightOwner_contain.cdc --network=testnet --args-json '[{"type": "String","value": "CopyrightOwner 001"}]'
```
#### Get CopyrightOwner Address by Name `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_copyrightOwner_address_by_name.cdc --network=testnet --args-json '[{"type": "String","value": "CopyrightOwner 001"}]'
```
#### Get Commonweal Names `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_commonweal_names.cdc --network=testnet
```
#### Get CommonwealCutPercentage by Name `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_commonwealCutPercentage_by_name.cdc --network=testnet --args-json '[{"type": "String","value": "Commonweal 001"}]'
```
#### Get Metaya Address `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_metaya_address.cdc --network=testnet
```
#### Get Metaya Market CutPercentage `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_metaya_marketCutPercentage.cdc --network=testnet
```
#### Get Store CutPercentages Amount `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_storeCutPercentages_amount.cdc --network=testnet
```
#### Get Store CutPercentages by Name `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_storeCutPercentage_by_name.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "String","value": "Metaya"}]'
```
#### Get Pack CutPercentages Amount `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_packCutPercentages_amount.cdc --network=testnet
```
#### Get Pack CutPercentages by Name `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_packCutPercentage_by_name.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "String","value": "Metaya"}]'
```
#### Get Market CutPercentages Amount `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_marketCutPercentages_amount.cdc --network=testnet
```
#### Get Market CutPercentages by Name `Script`
```
flow scripts execute ./scripts/metayaBeneficiaryCut/get_marketCutPercentage_by_name.cdc --network=testnet --args-json '[{"type": "UInt32","value": "1"},{"type": "String","value": "CopyrightOwner 001"}]'
```

# NFTStorefront Commands
The general-purpose contract is used in the Metaya market.

NFTStorefront contract is already deployed to testnet at [0x94b06cfca1d8a476](https://flow-view-source.com/testnet/account/0x94b06cfca1d8a476).

#### Setup Account `Transaction`
```
flow transactions send ./transactions/nftStorefront/setup_account.cdc --signer testnet-account --network=testnet
```
#### Sell Item by IUC `Transaction`
```
flow transactions send ./transactions/nftStorefront/sell_item_by_IUC.cdc --signer testnet-account --network=testnet --args-json '[{"type": "UInt64","value": "1"},{"type": "UFix64","value": "100.0"}]'
```
#### Remove Item `Transaction`
```
flow transactions send ./transactions/nftStorefront/remove_item.cdc --signer testnet-account --network=testnet --args-json '[{"type": "UInt64","value": "15320080"}]'
```
#### Cleanup Item `Transaction`
```
flow transactions send ./transactions/nftStorefront/cleanup_item.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt64","value": "23728806"},{"type": "Address","value": "0xd542317949eb00b6"}]'
```
#### Buy Item by Cash `Transaction`
```
flow transactions send ./transactions/nftStorefront/buy_item_by_cash.cdc --signer testnet-account-metaya --network=testnet --args-json '[{"type": "UInt64","value": "23728806"},{"type": "Address","value": "0xd542317949eb00b6"},{"type": "Address","value": "0x312588a458110069"}]'
```
#### Buy Item by IUC `Transaction`
```
flow transactions build ./transactions/nftStorefront/buy_item_by_IUC.cdc --network=testnet --args-json '[{"type": "UInt64","value": "23729111"},{"type": "Address","value": "0xd542317949eb00b6"}]' --authorizer testnet-account2 --authorizer testnet-account-metaya --proposer testnet-account2 --payer testnet-account-metaya --filter payload --save built.rlp

flow transactions sign ./built.rlp --signer testnet-account2 --network=testnet --filter payload --save signed.rlp

flow transactions sign ./signed.rlp --signer testnet-account-metaya --network=testnet --filter payload --save signed.rlp

flow transactions send-signed ./signed.rlp --network=testnet
```
#### Get Listing Ids `Script`
```
flow scripts execute ./scripts/nftStorefront/get_listing_ids.cdc --network=testnet --args-json '[{"type": "Address","value": "0xd542317949eb00b6"}]'
```
#### Get Listing Details `Script`
```
flow scripts execute ./scripts/nftStorefront/get_listing_details.cdc --network=testnet --args-json '[{"type": "Address","value": "0xd542317949eb00b6"},{"type": "UInt64","value": "23728806"}]'
```