import Metaya from "../../../contracts/Metaya.cdc"

// This transaction locks a set so that new plays can no longer be added to it

// Parameters:
//
// setID: the ID of the set to be locked

transaction(setID: UInt32) {

    // local variable for the admin resource
    let adminRef: &Metaya.Admin

    prepare(acct: AuthAccount) {

        // Borrow a reference to the admin resource
        self.adminRef = acct.borrow<&Metaya.Admin>(from: Metaya.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {

        // Borrow a reference to the Set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Lock the set permanently
        setRef.lock()
    }

    post {
        
        Metaya.isSetLocked(setID: setID)!:
            "Set did not lock"
    }
}