import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

transaction(playID: UInt32) {

    let adminRef: &MetayaBeneficiaryCut.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&MetayaBeneficiaryCut.Admin>(from: MetayaBeneficiaryCut.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {
        self.adminRef.setMarketCutPercentages(playID: playID, copyrightOwnerAndCutPercentage: nil)
    }
}