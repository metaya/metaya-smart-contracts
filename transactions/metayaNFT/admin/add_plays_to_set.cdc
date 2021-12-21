import Metaya from "../../../contracts/Metaya.cdc"

// This transaction adds multiple plays to a set

// Parameters:
//
// setID: the ID of the set to which multiple plays are added
// plays: an array of play IDs being added to the set

transaction(setID: UInt32, plays: [UInt32]) {

    // Local variable for the Metaya Admin object
    let adminRef: &Metaya.Admin

    prepare(acct: AuthAccount) {

        // Borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&Metaya.Admin>(from: Metaya.AdminStoragePath)!
    }

    execute {

        // Borrow a reference to the set to be added to
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Add the specified play IDs
        setRef.addPlays(playIDs: plays)
    }
}