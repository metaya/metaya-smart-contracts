import FungibleToken from "../../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"
import Metaya from "../../../contracts/Metaya.cdc"
import MetayaShardedCollection from "../../../contracts/MetayaShardedCollection.cdc"
import FlowStorageFees from "../../../contracts/FlowStorageFees.cdc"
import FlowToken from "../../../contracts/FlowToken.cdc"

// This transaction is what Metaya uses to send the moments in a "pack" to
// a user's collection

// Parameters:
//
// recipientAddr: the Flow address of the account receiving a pack of moments
// momentsIDs: an array of moment IDs to be withdrawn from the owner's moment collection

transaction(recipientAddr: Address, momentIDs: [UInt64]) {

    prepare(acct: AuthAccount) {

        // Used to determine whether signer needs more storage
        fun returnFlowFromStorage(_ storage: UInt64): UFix64 {

            // Safe convert UInt64 to UFix64 (without overflow)
            let f = UFix64(storage % 100000000 as UInt64) * 0.00000001 as UFix64 + UFix64(storage / 100000000 as UInt64)

            // Decimal point correction. Megabytes to bytes have a conversion of 10^-6 while UFix64 minimum value is 10^-8
            let storageMb = f * 100.0 as UFix64
            let storage = FlowStorageFees.storageCapacityToFlow(storageMb)
            return storage
        }

        // Get the recipient's public account object
        let recipient = getAccount(recipientAddr)

        // Borrow a reference to the recipient's moment collection
        let receiverRef = recipient.getCapability(Metaya.CollectionPublicPath)
            .borrow<&{Metaya.MomentCollectionPublic}>()
            ?? panic("Could not borrow reference to receiver's collection")

        // Borrow a reference to the owner's moment collection
        if let collection = acct.borrow<&MetayaShardedCollection.ShardedCollection>(from: MetayaShardedCollection.ShardedCollectionStoragePath) {

            receiverRef.batchDeposit(tokens: <-collection.batchWithdraw(ids: momentIDs))

        } else {

            let collection = acct.borrow<&Metaya.Collection>(from: Metaya.CollectionStoragePath)!

            // Deposit the pack of moments to the recipient's collection
            receiverRef.batchDeposit(tokens: <-collection.batchWithdraw(ids: momentIDs))

        }

        // Determine Storage Used by user and Total Capacity in their account
        var storageUsed = returnFlowFromStorage(recipient.storageUsed)
        var storageTotal = returnFlowFromStorage(recipient.storageCapacity)

        // If user has used more than their total capacity, increase total capacity
        if (storageUsed > storageTotal) {

            let difference = storageUsed - storageTotal

            // Withdraw storage fee
            let vaultRef = acct.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
                ?? panic("Failed to borrow reference to sender vault")
            let sentVault <- vaultRef.withdraw(amount: difference)

            // Deposit storage fee to recipient
            let receiver = recipient.getCapability(/public/flowTokenReceiver)
                .borrow<&{FungibleToken.Receiver}>()
                    ?? panic("Failed to borrow reference to recipient vault")
            receiver.deposit(from: <-sentVault)
        }
    }
}