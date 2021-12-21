import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"
import Metaya from "../../contracts/Metaya.cdc"
import NFTStorefront from "../../contracts/NFTStorefront.cdc"
import FlowStorageFees from "../../contracts/FlowStorageFees.cdc"
import FlowToken from "../../contracts/FlowToken.cdc"

transaction(listingResourceID: UInt64, storefrontAddress: Address) {

    let storefront: &NFTStorefront.Storefront{NFTStorefront.StorefrontPublic}
    let listing: &NFTStorefront.Listing{NFTStorefront.ListingPublic}
    let paymentVault: @FungibleToken.Vault
    let buyerAccount: PublicAccount
    let metayaNFTCollection: &Metaya.Collection{NonFungibleToken.Receiver}
    let vaultRef: &FlowToken.Vault

    prepare(acct: AuthAccount, admin: AuthAccount) {
        self.storefront = getAccount(storefrontAddress)
            .getCapability<&NFTStorefront.Storefront{NFTStorefront.StorefrontPublic}>(
                NFTStorefront.StorefrontPublicPath
            )!
            .borrow()
            ?? panic("Could not borrow Storefront from provided address")
        self.listing = self.storefront.borrowListing(listingResourceID: listingResourceID)
                    ?? panic("No offer with that ID in Storefront")
        let price = self.listing.getDetails().salePrice

        self.buyerAccount = getAccount(acct.address)

        let metayaUtilityCoinVault = acct.borrow<&MetayaUtilityCoin.Vault>(from: MetayaUtilityCoin.VaultStoragePath)
            ?? panic("Cannot borrow MetayaUtilityCoin vault from account storage")
        self.paymentVault <- metayaUtilityCoinVault.withdraw(amount: price)

        self.metayaNFTCollection = acct.borrow<&Metaya.Collection{NonFungibleToken.Receiver}>(from: Metaya.CollectionStoragePath)
            ?? panic("Cannot borrow Metaya NFT collection receiver from account")

        // Borrow a reference to the admin flow token vault
        self.vaultRef = admin.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
            ?? panic("Failed to borrow reference to admin vault")
    }

    execute {
        let item <- self.listing.purchase(
            payment: <- self.paymentVault
        )

        self.metayaNFTCollection.deposit(token: <-item)

        // Be kind and recycle
        self.storefront.cleanup(listingResourceID: listingResourceID)

        // Used to determine whether buyerAccount needs more storage
        fun returnFlowFromStorage(_ storage: UInt64): UFix64 {
            // Safe convert UInt64 to UFix64 (without overflow)
            let f = UFix64(storage % 100000000 as UInt64) * 0.00000001 as UFix64 + UFix64(storage / 100000000 as UInt64)
            // Decimal point correction. Megabytes to bytes have a conversion of 10^-6 while UFix64 minimum value is 10^-8
            let storageMb = f * 100.0 as UFix64
            let storage = FlowStorageFees.storageCapacityToFlow(storageMb)
            return storage
        }

        // Determine Storage Used by user and Total Capacity in their account
        var storageUsed = returnFlowFromStorage(self.buyerAccount.storageUsed)
        var storageTotal = returnFlowFromStorage(self.buyerAccount.storageCapacity)

        // If user has used more than their total capacity, increase total capacity
        if (storageUsed > storageTotal) {
            let difference = storageUsed - storageTotal
            // Withdraw storage fee
            let sentVault <- self.vaultRef.withdraw(amount: difference)

            // Deposit storage fee to buyerAccount
            let receiver = self.buyerAccount.getCapability(/public/flowTokenReceiver)
                .borrow<&{FungibleToken.Receiver}>()
                    ?? panic("Failed to borrow reference to buyerAccount vault")
            receiver.deposit(from: <-sentVault)
        }
    }
}