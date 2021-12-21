import Metaya from "../../../contracts/Metaya.cdc"

// This script reads the public nextPlayID from the Metaya contract and
// returns that number to the caller

// Returns: UInt32
// the nextPlayID field in Metaya contract

pub fun main(): UInt32 {

    log(Metaya.nextPlayID)

    return Metaya.nextPlayID
}