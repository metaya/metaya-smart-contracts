import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import SoundlinksDID from "../../contracts/SoundlinksDID.cdc"
import FungibleToken from "../../contracts/FungibleToken.cdc"
import FlowToken from "../../contracts/FlowToken.cdc"

transaction(purchaseAmount: UInt32, hashs: [String], purchaseUnitPrice: UFix64) {

    let didAdmin: &SoundlinksDID.Admin
    let didReceiver: &{NonFungibleToken.CollectionPublic}
    let flowPayer: &FlowToken.Vault
    let flowReceiver: &{FungibleToken.Receiver}

    prepare(soundlinksAdmin: AuthAccount, signer: AuthAccount) {

        if signer.borrow<&SoundlinksDID.Collection>(from: SoundlinksDID.CollectionStoragePath) == nil {

            signer.save(
                <-SoundlinksDID.createEmptyCollection(),
                to: SoundlinksDID.CollectionStoragePath
            )

            signer.link<&SoundlinksDID.Collection{NonFungibleToken.CollectionPublic, SoundlinksDID.SoundlinksDIDCollectionPublic}>(
                SoundlinksDID.CollectionPublicPath,
                target: SoundlinksDID.CollectionStoragePath
            )
        }

        self.didAdmin = soundlinksAdmin
            .borrow<&SoundlinksDID.Admin>(from: SoundlinksDID.AdminStoragePath)
            ?? panic("soundlinksAdmin is not the Soundlinks DID admin.")

        self.didReceiver = signer
            .getCapability(SoundlinksDID.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not borrow receiver reference to the recipient's DID Collection.")

        self.flowPayer = signer
            .borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
            ?? panic("Failed to borrow reference to signer's Flow Vault.")

        self.flowReceiver = soundlinksAdmin
            .getCapability(/public/flowTokenReceiver)!
            .borrow<&{FungibleToken.Receiver}>()
            ?? panic("Could not borrow receiver reference to the recipient's Flow Vault.")
    }

    pre {

        UInt32(hashs.length) == purchaseAmount: "The amount of hashs should be the same as the purchaseAmount."
    }

    execute {

        let amount = UFix64(purchaseAmount) * purchaseUnitPrice
        let sentVault <- self.flowPayer.withdraw(amount: amount)
        self.flowReceiver.deposit(from: <- sentVault)

        self.didAdmin.mintDIDs(recipient: self.didReceiver, hashs: hashs)
    }
}