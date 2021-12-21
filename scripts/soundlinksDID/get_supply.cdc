import SoundlinksDID from "../../contracts/SoundlinksDID.cdc"

pub fun main(): UInt64 {

    let supply = SoundlinksDID.totalSupply

    log(supply)

    return supply
}