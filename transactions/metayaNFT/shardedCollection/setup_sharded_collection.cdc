import Metaya from "../../../contracts/Metaya.cdc"
import MetayaShardedCollection from "../../../contracts/MetayaShardedCollection.cdc"

// This transaction creates and stores an empty moment collection
// and creates a public capability for it.
// Moments are split into a number of buckets
// This makes storage more efficient and performant

// Parameters
//
// numBuckets: The number of buckets to split Moments into

transaction(numBuckets: UInt64) {

    prepare(acct: AuthAccount) {

        if acct.borrow<&MetayaShardedCollection.ShardedCollection>(from: MetayaShardedCollection.ShardedCollectionStoragePath) == nil {

            let collection <- MetayaShardedCollection.createEmptyCollection(numBuckets: numBuckets)

            // Put a new Collection in storage
            acct.save(<-collection, to: MetayaShardedCollection.ShardedCollectionStoragePath)

            // Create a public capability for the collection
            if acct.link<&{Metaya.MomentCollectionPublic}>(Metaya.CollectionPublicPath, target: MetayaShardedCollection.ShardedCollectionStoragePath) == nil {
                acct.unlink(Metaya.CollectionPublicPath)
            }

            acct.link<&{Metaya.MomentCollectionPublic}>(Metaya.CollectionPublicPath, target: MetayaShardedCollection.ShardedCollectionStoragePath)
        } else {

            panic("Sharded Collection already exists!")
        }
    }
}