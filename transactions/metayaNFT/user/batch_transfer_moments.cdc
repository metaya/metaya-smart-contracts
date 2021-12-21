import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"
import Metaya from "../../../contracts/Metaya.cdc"
import NFTStorefront from "../../../contracts/NFTStorefront.cdc"

// This transaction transfers a number of moments to a recipient

// Parameters
//
// recipientAddress: the Flow address who will receive the NFTs
// momentIDs: an array of moment IDs of NFTs that recipient will receive

transaction(recipient: Address, momentIDs: [UInt64]) {

    let transferTokens: @NonFungibleToken.Collection
    
    prepare(acct: AuthAccount) {

        self.transferTokens <- acct.borrow<&Metaya.Collection>(from: Metaya.CollectionStoragePath)!.batchWithdraw(ids: momentIDs)

        if let saleRef = acct.borrow<&NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath) {

            for id in saleRef.getListingIDs() {

                let listing = saleRef.borrowListing(listingResourceID: id)!

                if listing.getDetails().nftType == Type<@Metaya.NFT>() {

                    for withdrawID in momentIDs {

                        if listing.getDetails().nftID == withdrawID {

                            saleRef.removeListing(listingResourceID: id)
                            break
                        }
                    }
                }
            }
        }
    }

    execute {

        // Get the recipient's public account object
        let recipient = getAccount(recipient)

        // Get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(Metaya.CollectionPublicPath).borrow<&{Metaya.MomentCollectionPublic}>()
            ?? panic("Could not borrow a reference to the recipients moment receiver")

        // Deposit the NFT in the receivers collection
        receiverRef.batchDeposit(tokens: <-self.transferTokens)
    }
}