import MetayaBeneficiaryCut from "../../contracts/MetayaBeneficiaryCut.cdc"

pub fun main(name: String): Bool {

    return MetayaBeneficiaryCut.getAllCopyrightOwnerNames().contains(name)

}