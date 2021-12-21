import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import MetayaUtilityCoin from "../../contracts/MetayaUtilityCoin.cdc"
import Metaya from "../../contracts/Metaya.cdc"
import NFTStorefront from "../../contracts/NFTStorefront.cdc"
import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

transaction(saleItemID: UInt64, saleItemPrice: UFix64) {

    let metayaUtilityCoinReceiver: Capability<&MetayaUtilityCoin.Vault{FungibleToken.Receiver}>
    let metayaNFTProvider: Capability<&Metaya.Collection{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>
    let metayaNFTcollectionRef: &Metaya.Collection
    let storefront: &NFTStorefront.Storefront

    prepare(acct: AuthAccount) {

        // We need a provider capability, but one is not provided by default so we create one if needed.
        let metayaNFTCollectionProviderPrivatePath = /private/metayaNFTCollectionProviderForNFTStorefront

        self.metayaUtilityCoinReceiver = acct.getCapability<&MetayaUtilityCoin.Vault{FungibleToken.Receiver}>(MetayaUtilityCoin.ReceiverPublicPath)!

        assert(self.metayaUtilityCoinReceiver.borrow() != nil, message: "Missing or mis-typed MetayaUtilityCoin receiver")

        if !acct.getCapability<&Metaya.Collection{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(metayaNFTCollectionProviderPrivatePath)!.check() {
            acct.link<&Metaya.Collection{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(metayaNFTCollectionProviderPrivatePath, target: Metaya.CollectionStoragePath)
        }

        self.metayaNFTProvider = acct.getCapability<&Metaya.Collection{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(metayaNFTCollectionProviderPrivatePath)!
        assert(self.metayaNFTProvider.borrow() != nil, message: "Missing or mis-typed Metaya.Collection provider")

        self.metayaNFTcollectionRef = acct.borrow<&Metaya.Collection>(from: Metaya.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the stored Moment collection")

        self.storefront = acct.borrow<&NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath)
            ?? panic("Missing or mis-typed NFTStorefront Storefront")
    }

    execute {

        var saleCuts: [NFTStorefront.SaleCut] = []
        var sellerCutAmount = saleItemPrice

        let nftRef = self.metayaNFTcollectionRef.borrowMoment(id: saleItemID)
            ?? panic("Could not borrow a reference to the NFT")
        let playID = nftRef.data.playID

        // Metaya Market Cut
        let metayaMarketCutPercentage = MetayaBeneficiaryCut.metayaMarketCutPercentage
        let metayaMarketCutAmount = saleItemPrice * metayaMarketCutPercentage
        saleCuts.append(NFTStorefront.SaleCut(
                receiver: MetayaBeneficiaryCut.metayaCapability,
                amount: metayaMarketCutAmount
        ))
        sellerCutAmount = sellerCutAmount - metayaMarketCutAmount

        // Copyright owners Market Cut
        for name in MetayaBeneficiaryCut.getMarketCopyrightOwnerNames(playID: playID)! {
            let copyrightOwnerCap = MetayaBeneficiaryCut.getCopyrightOwnerCapability(name: name)
                ?? panic("Cannot find the copyright owner by the name.")
            let copyrightOwnerCutPercentage = MetayaBeneficiaryCut.getMarketCutPercentage(playID: playID, name: name)
                ?? panic("Cannot find the copyright owner cutPercentage by the name.")
            let copyrightOwnerCutAmount = saleItemPrice * copyrightOwnerCutPercentage
            saleCuts.append(NFTStorefront.SaleCut(
                receiver: copyrightOwnerCap as Capability<&{FungibleToken.Receiver}>,
                amount: copyrightOwnerCutAmount
            ))
            sellerCutAmount = sellerCutAmount - copyrightOwnerCutAmount
        }

        // Seller Cut
        let sellerCut = NFTStorefront.SaleCut(
            receiver: self.metayaUtilityCoinReceiver,
            amount: sellerCutAmount
        )
        saleCuts.insert(at:0 ,sellerCut)

        self.storefront.createListing(
            nftProviderCapability: self.metayaNFTProvider,
            nftType: Type<@Metaya.NFT>(),
            nftID: saleItemID,
            salePaymentVaultType: Type<@MetayaUtilityCoin.Vault>(),
            saleCuts: saleCuts
        )
    }
}