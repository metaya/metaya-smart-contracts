import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import SoundlinksDID from "../../contracts/SoundlinksDID.cdc"

transaction(recipient: Address, hashs: [String]) {

    let didAdmin: &SoundlinksDID.Admin

    prepare(signer: AuthAccount) {

        self.didAdmin = signer
            .borrow<&SoundlinksDID.Admin>(from: SoundlinksDID.AdminStoragePath)
            ?? panic("Signer is not the Soundlinks DID admin")
    }

    execute {

        let didReceiver = getAccount(recipient)
            .getCapability(SoundlinksDID.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Unable to borrow receiver reference")

        self.didAdmin.mintDIDs(recipient: didReceiver, hashs: hashs)
    }
}