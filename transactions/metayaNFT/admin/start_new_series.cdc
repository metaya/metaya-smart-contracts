import Metaya from "../../../contracts/Metaya.cdc"

// This transaction is for an Admin to start a new Metaya series

transaction {

    // Local variable for the Metaya Admin object
    let adminRef: &Metaya.Admin
    let currentSeries: UInt32

    prepare(acct: AuthAccount) {

        // Borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&Metaya.Admin>(from: Metaya.AdminStoragePath)
            ?? panic("No admin resource in storage")

        self.currentSeries = Metaya.currentSeries
    }

    execute {

        // Increment the series number
        self.adminRef.startNewSeries()
    }

    post {

        Metaya.currentSeries == self.currentSeries + 1 as UInt32:
            "new series not started"
    }
}