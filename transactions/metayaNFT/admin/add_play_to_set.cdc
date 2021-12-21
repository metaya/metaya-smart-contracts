import Metaya from "../../../contracts/Metaya.cdc"

// This transaction is how a Metaya admin adds a created play to a set

// Parameters:
//
// setID: the ID of the set to which a created play is added
// playID: the ID of the play being added

transaction(setID: UInt32, playID: UInt32) {

    // Local variable for the Metaya Admin object
    let adminRef: &Metaya.Admin

    prepare(acct: AuthAccount) {

    // Borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&Metaya.Admin>(from: Metaya.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Admin resource")
    }

    execute {

        // Borrow a reference to the set to be added to
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Add the specified play ID
        setRef.addPlay(playID: playID)
    }

    post {
        Metaya.getPlaysInSet(setID: setID)!.contains(playID): 
            "set does not contain playID"
    }
}