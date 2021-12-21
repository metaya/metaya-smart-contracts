import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(saleID: UInt32, name: String): UFix64? {

    return MetayaBeneficiaryCut.getStoreCutPercentage(saleID: saleID, name: name)

}