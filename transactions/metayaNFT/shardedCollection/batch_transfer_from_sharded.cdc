import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"
import Metaya from "../../../contracts/Metaya.cdc"
import MetayaShardedCollection from "../../../contracts/MetayaShardedCollection.cdc"

// This transaction deposits a number of NFTs to a recipient

// Parameters
//
// recipient: the Flow address who will receive the NFTs
// momentIDs: an array of moment IDs of NFTs that recipient will receive

transaction(recipient: Address, momentIDs: [UInt64]) {

    let transferTokens: @NonFungibleToken.Collection

    prepare(acct: AuthAccount) {

        self.transferTokens <- acct.borrow<&MetayaShardedCollection.ShardedCollection>(from: MetayaShardedCollection.ShardedCollectionStoragePath)!.batchWithdraw(ids: momentIDs)
    }

    execute {

        // Get the recipient's public account object
        let recipient = getAccount(recipient)

        // Get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(Metaya.CollectionPublicPath).borrow<&{Metaya.MomentCollectionPublic}>()!

        // Deposit the NFT in the receivers collection
        receiverRef.batchDeposit(tokens: <-self.transferTokens)
    }
}