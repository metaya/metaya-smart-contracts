import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

transaction(saleID: UInt32) {

    let adminRef: &MetayaBeneficiaryCut.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&MetayaBeneficiaryCut.Admin>(from: MetayaBeneficiaryCut.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {
        self.adminRef.setStoreCutPercentages(saleID: saleID, copyrightOwnerAndCutPercentage: nil)
    }
}