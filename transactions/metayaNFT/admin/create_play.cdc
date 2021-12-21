import Metaya from "../../../contracts/Metaya.cdc"

// This transaction creates a new play struct
// and stores it in the Metaya smart contract
// We currently stringify the metadata and insert it into the
// transaction string

// Parameters:
//
// metadata: A dictionary of all the play metadata associated

transaction(metadata: {String: String}) {

    // Local variable for the Metaya Admin object
    let adminRef: &Metaya.Admin
    let currPlayID: UInt32

    prepare(acct: AuthAccount) {

        // Borrow a reference to the admin resource
        self.currPlayID = Metaya.nextPlayID
        self.adminRef = acct.borrow<&Metaya.Admin>(from: Metaya.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {

        // Create a play with the specified metadata
        self.adminRef.createPlay(metadata: metadata)
    }

    post {

        Metaya.getPlayMetaData(playID: self.currPlayID) != nil:
            "playID doesnt exist"
    }
}