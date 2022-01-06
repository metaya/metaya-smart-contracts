/**
    Description: Central Smart Contract for Metaya Admin Receiver
    This contract defines a function that takes a Metaya Admin
    object and stores it in the storage of the contract account
    so it can be used.

    Copyright 2021 Metaya.io
    SPDX-License-Identifier: Apache-2.0
**/

import Metaya from 0x8b935cd43003d4b2
import MetayaShardedCollection from 0x8b935cd43003d4b2

pub contract MetayaAdminReceiver {

    /// storeAdmin takes a Metaya Admin resource and 
    /// saves it to the account storage of the account
    /// where the contract is deployed
    pub fun storeAdmin(newAdmin: @Metaya.Admin) {
        self.account.save(<-newAdmin, to: Metaya.AdminStoragePath)
    }

    init() {
        // Save a copy of the sharded Moment Collection to the account storage
        if self.account.borrow<&MetayaShardedCollection.ShardedCollection>(from: MetayaShardedCollection.ShardedCollectionStoragePath) == nil {
            let collection <- MetayaShardedCollection.createEmptyCollection(numBuckets: 32)
            // Put a new Collection in storage
            self.account.save(<-collection, to: MetayaShardedCollection.ShardedCollectionStoragePath)

            self.account.link<&{Metaya.MomentCollectionPublic}>(Metaya.CollectionPublicPath, target: MetayaShardedCollection.ShardedCollectionStoragePath)
        }
    }
}