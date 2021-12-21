import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(): Int {

    return MetayaBeneficiaryCut.getSaleIDsInStoreCutPercentages().length

}