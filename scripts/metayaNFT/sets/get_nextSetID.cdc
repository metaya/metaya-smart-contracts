import Metaya from "../../../contracts/Metaya.cdc"

// This script reads the next Set ID from the Metaya contract and
// returns that number to the caller

// Returns: UInt32
// Value of nextSetID field in Metaya contract

pub fun main(): UInt32 {

    log(Metaya.nextSetID)

    return Metaya.nextSetID
}