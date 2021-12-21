import Metaya from "../../../contracts/Metaya.cdc"

// This transaction is for retiring all plays from a set, which
// makes it so that moments can no longer be minted
// from all the editions with that set

// Parameters:
//
// setID: the ID of the set to be retired entirely

transaction(setID: UInt32) {

    // Local variable for the admin reference
    let adminRef: &Metaya.Admin

    prepare(acct: AuthAccount) {

        // Borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&Metaya.Admin>(from: Metaya.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {

        // Borrow a reference to the specified set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Retire all the plays
        setRef.retireAll()
    }
}