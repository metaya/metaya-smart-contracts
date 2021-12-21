import Metaya from "../../../contracts/Metaya.cdc"

// This transaction sets up an account to use Metaya
// by storing an empty moment collection and creating
// a public capability for it

transaction {

    prepare(acct: AuthAccount) {

        // First, check to see if a moment collection already exists
        if acct.borrow<&Metaya.Collection>(from: Metaya.CollectionStoragePath) == nil {

            // Create a new Metaya Collection
            let collection <- Metaya.createEmptyCollection() as! @Metaya.Collection

            // Put the new Collection in storage
            acct.save(<-collection, to: Metaya.CollectionStoragePath)

            // create a public capability for the collection
            acct.link<&{Metaya.MomentCollectionPublic}>(Metaya.CollectionPublicPath, target: Metaya.CollectionStoragePath)
        }
    }
}