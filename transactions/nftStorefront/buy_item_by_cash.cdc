import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"
import Metaya from "../../contracts/Metaya.cdc"
import NFTStorefront from "../../contracts/NFTStorefront.cdc"
import FlowStorageFees from "../../contracts/FlowStorageFees.cdc"
import FlowToken from "../../contracts/FlowToken.cdc"

transaction(listingResourceID: UInt64, storefrontAddress: Address, buyerAddress: Address) {

    let storefront: &NFTStorefront.Storefront{NFTStorefront.StorefrontPublic}
    let listing: &NFTStorefront.Listing{NFTStorefront.ListingPublic}
    let price: UFix64
    let ducRef: &MetayaUtilityCoin.Administrator
    let buyerAccount: PublicAccount
    let metayaNFTCollection: &{Metaya.MomentCollectionPublic}
    let vaultRef: &FlowToken.Vault

    prepare(admin: AuthAccount) {
        self.storefront = getAccount(storefrontAddress)
            .getCapability<&NFTStorefront.Storefront{NFTStorefront.StorefrontPublic}>(
                NFTStorefront.StorefrontPublicPath
            )!
            .borrow()
            ?? panic("Could not borrow Storefront from provided address")
        self.listing = self.storefront.borrowListing(listingResourceID: listingResourceID)
                    ?? panic("No offer with that ID in Storefront")
        self.price = self.listing.getDetails().salePrice

        self.ducRef = admin
            .borrow<&MetayaUtilityCoin.Administrator>(from: MetayaUtilityCoin.AdminStoragePath)
            ?? panic("Signer is not the token admin")

        self.buyerAccount = getAccount(buyerAddress)
        self.metayaNFTCollection = self.buyerAccount.getCapability(Metaya.CollectionPublicPath).borrow<&{Metaya.MomentCollectionPublic}>()
            ?? panic("Could not borrow a reference to the moment collection")

        // Borrow a reference to the admin flow token vault
        self.vaultRef = admin.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
            ?? panic("Failed to borrow reference to admin vault")
    }

    execute {
        let minter <- self.ducRef.createNewMinter(allowedAmount: self.price)
        let paymentVault <- minter.mintTokens(amount: self.price) as! @MetayaUtilityCoin.Vault
        destroy minter

        let item <- self.listing.purchase(
            payment: <- paymentVault
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