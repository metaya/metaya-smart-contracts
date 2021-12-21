import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

transaction(cutPercentage: UFix64) {

    let adminRef: &MetayaBeneficiaryCut.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&MetayaBeneficiaryCut.Admin>(from: MetayaBeneficiaryCut.AdminStoragePath)
            ?? panic("No admin resource in storage")
    }

    execute {
        self.adminRef.setMetayaMarketCutPercentage(cutPercentage: cutPercentage)
    }
}