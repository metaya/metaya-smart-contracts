import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(playID: UInt32, name: String): UFix64? {

    return MetayaBeneficiaryCut.getMarketCutPercentage(playID: playID, name: name)

}