import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

transaction(packID: UInt32, copyrightOwnerAndCutPercentage: {String: UFix64}) {

    let adminRef: &MetayaBeneficiaryCut.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&MetayaBeneficiaryCut.Admin>(from: MetayaBeneficiaryCut.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {
        self.adminRef.setPackCutPercentages(packID: packID, copyrightOwnerAndCutPercentage: copyrightOwnerAndCutPercentage)
    }
}