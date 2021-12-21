import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(packID: UInt32, name: String): UFix64? {

    return MetayaBeneficiaryCut.getPackCutPercentage(packID: packID, name: name)

}