import Metaya from "../../../contracts/Metaya.cdc"
import SoundlinksDID from "../../../contracts/SoundlinksDID.cdc"

// This transaction mints multiple moments
// from a single set/play combination (otherwise known as edition)

// Parameters:
//
// setID: the ID of the set to be minted from
// playID: the ID of the Play from which the Moments are minted
// quantity: the quantity of Moments to be minted
// recipientAddr: the Flow address of the account receiving the collection of minted moments

transaction(setID: UInt32, playID: UInt32, quantity: UInt32, recipientAddr: Address) {

    // Local variable for the Metaya Admin object
    let adminRef: &Metaya.Admin
    let soundlinksDIDCollection: &SoundlinksDID.Collection

    prepare(acct: AuthAccount) {

        // Borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&Metaya.Admin>(from: Metaya.AdminStoragePath)!
        self.soundlinksDIDCollection = acct.borrow<&SoundlinksDID.Collection>(from: SoundlinksDID.CollectionStoragePath)!
    }

    execute {

        // borrow a reference to the set to be minted from
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Get SOUNDLINKS DIDs
        let ids = self.soundlinksDIDCollection.getIDsByAmount(amount: quantity)
        let transferDIDCollection <- self.soundlinksDIDCollection.batchWithdraw(ids: ids) as! @SoundlinksDID.Collection

        // Mint all the new NFTs
        let collection <- setRef.batchMintMoment(playID: playID, quantity: quantity, soundlinksDIDCollection: <-transferDIDCollection)

        // Get the account object for the recipient of the minted tokens
        let recipient = getAccount(recipientAddr)

        // Get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(Metaya.CollectionPublicPath).borrow<&{Metaya.MomentCollectionPublic}>()
            ?? panic("Cannot borrow a reference to the recipient's collection")

        // Deposit the NFT in the receivers collection
        receiverRef.batchDeposit(tokens: <-collection)
    }
}