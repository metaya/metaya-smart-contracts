import SoundlinksDID from "../../contracts/SoundlinksDID.cdc"

pub fun main(address: Address): UInt32 {

    let account = getAccount(address)

    let collectionRef = account.getCapability(SoundlinksDID.CollectionPublicPath)!
        .borrow<&SoundlinksDID.Collection{SoundlinksDID.SoundlinksDIDCollectionPublic}>()
        ?? panic("Could not borrow the reference to the Collection")

    return UInt32(collectionRef.getIDs().length)
}