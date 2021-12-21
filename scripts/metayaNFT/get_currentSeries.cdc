import Metaya from "../../contracts/Metaya.cdc"

// This script reads the current series from the Metaya contract and
// returns that number to the caller

// Returns: UInt32
// currentSeries field in Metaya contract

pub fun main(): UInt32 {

    return Metaya.currentSeries
}