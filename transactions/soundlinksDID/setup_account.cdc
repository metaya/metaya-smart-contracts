import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import SoundlinksDID from "../../contracts/SoundlinksDID.cdc"

transaction {

    prepare(signer: AuthAccount) {

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
    }
}