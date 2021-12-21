import Metaya from "../../../contracts/Metaya.cdc"
import MetayaAdminReceiver from "../../../contracts/MetayaAdminReceiver.cdc"

// This transaction takes a Metaya Admin resource and
// saves it to the account storage of the account
// where the contract is deployed

transaction {

    // Local variable for the Metaya Admin object
    let adminRef: @Metaya.Admin

    prepare(acct: AuthAccount) {

        self.adminRef <- acct.load<@Metaya.Admin>(from: Metaya.AdminStoragePath)
            ?? panic("No Metaya admin in storage")

    }

    execute {

        MetayaAdminReceiver.storeAdmin(newAdmin: <-self.adminRef)

    }
}