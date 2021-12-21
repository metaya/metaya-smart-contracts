import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(name: String): UFix64? {

    return MetayaBeneficiaryCut.getCommonwealCutPercentage(name: name)

}