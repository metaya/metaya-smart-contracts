import Metaya from "../../../contracts/Metaya.cdc"

// This transaction is for the admin to create a new set resource
// and store it in the Metaya smart contract

// Parameters:
//
// setName: the name of a new Set to be created

transaction(setName: String) {

    // Local variable for the Metaya Admin object
    let adminRef: &Metaya.Admin
    let currSetID: UInt32

    prepare(acct: AuthAccount) {

        // Borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&Metaya.Admin>(from: Metaya.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Admin resource")
        self.currSetID = Metaya.nextSetID
    }

    execute {

        // Create a set with the specified name
        self.adminRef.createSet(name: setName)
    }

    post {

        Metaya.getSetName(setID: self.currSetID) == setName:
          "Could not find the specified set"
    }
}