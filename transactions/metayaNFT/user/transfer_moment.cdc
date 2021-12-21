import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"
import Metaya from "../../../contracts/Metaya.cdc"
import NFTStorefront from "../../../contracts/NFTStorefront.cdc"

// This transaction transfers a moment to a recipient
// and cancels the sale in the collection if it exists

// Parameters:
//
// recipient: The Flow address of the account to receive the moment.
// withdrawID: The id of the moment to be transferred

transaction(recipient: Address, withdrawID: UInt64) {

    // Local variable for storing the transferred token
    let transferToken: @NonFungibleToken.NFT

    prepare(acct: AuthAccount) {

        // Borrow a reference to the owner's collection
        let collectionRef = acct.borrow<&Metaya.Collection>(from: Metaya.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the stored Moment collection")

        // Withdraw the NFT
        self.transferToken <- collectionRef.withdraw(withdrawID: withdrawID)

        if let saleRef = acct.borrow<&NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath) {

            for id in saleRef.getListingIDs() {

                let listing = saleRef.borrowListing(listingResourceID: id)!

                if listing.getDetails().nftType == Type<@Metaya.NFT>() {

                    if listing.getDetails().nftID == withdrawID {

                        saleRef.removeListing(listingResourceID: id)
                    }
                }
            }
        }
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