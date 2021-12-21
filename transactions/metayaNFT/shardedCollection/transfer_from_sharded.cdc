import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"
import Metaya from "../../../contracts/Metaya.cdc"
import MetayaShardedCollection from "../../../contracts/MetayaShardedCollection.cdc"

// This transaction deposits an NFT to a recipient

// Parameters
//
// recipient: the Flow address who will receive the NFT
// momentID: moment ID of NFT that recipient will receive

transaction(recipient: Address, momentID: UInt64) {

    let transferToken: @NonFungibleToken.NFT

    prepare(acct: AuthAccount) {

        self.transferToken <- acct.borrow<&MetayaShardedCollection.ShardedCollection>(from: MetayaShardedCollection.ShardedCollectionStoragePath)!.withdraw(withdrawID: momentID)
    }

    execute {

        // Get the recipient's public account object
        let recipient = getAccount(recipient)

        // Get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(Metaya.CollectionPublicPath).borrow<&{Metaya.MomentCollectionPublic}>()!

        // Deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-self.transferToken)
    }
}