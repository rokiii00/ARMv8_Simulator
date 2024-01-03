/*
CS3339- Computer Architecture. (MW 2-3:20pm)
Project 2: ARM V8 decoder and simulator.

Team Members: Matthew Lee, Hunter Savage-Pierce, Alan Solis.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

var if_break bool
var PC_Address int
var DataStartIndex int
var arrIndex int
var RegArray [32]int
var afterBreakData []int32

func main() {

	var InputFileName *string
	var OutputFileName *string
	var SimulationFileName string

	//Defines flags for arguments from the console - Alan
	InputFileName = flag.String("i", "", "Gets the input file name")
	OutputFileName = flag.String("o", "", "Gets the output file name")

	//Parses arguments from console - Alan
	flag.Parse()

	//create simulation output file name
	SimulationFileName = *OutputFileName

	//append different endings to output files
	*OutputFileName += "_dis.txt"
	SimulationFileName += "_sim.txt"

	//If there are too few or too many arguments that are parsed, the program will exit - Alan
	if flag.NArg() != 0 {
		os.Exit(200)
	}

	//opens the binary files.-Matt
	Contents, err := os.Open(*InputFileName)
	//closes the file-Matt
	defer Contents.Close()

	// error opening file (Base Case).- Matt.
	if err != nil {
		fmt.Println("ERROR opening file: ", *InputFileName)
		os.Exit(-1)
	}

	//creates output file - Matt.
	Storage, err := os.Create(*OutputFileName)
	//closes output file at end of program - Matt
	defer Storage.Close()

	//if program cannot create output file, then the system prints this - Matt.
	if err != nil {
		fmt.Println("Error creating file: ", *OutputFileName)
		panic(err)
	}

	//creates simulation file- Alan
	Simulator, err := os.Create(SimulationFileName)
	//closes simulation file at end of program - Alan
	defer Simulator.Close()

	//if program cannot create simulation file, then the system prints this - Alan.
	if err != nil {
		fmt.Println("Error creating file: ", SimulationFileName)
		panic(err)
	}

	//establishes a new file that reads all inputs line by line.
	Scanner := bufio.NewScanner(Contents)
	var writeToOutput string
	var binaryString string
	var opcodeString string
	var TotalOutput []string

	// Modified Matts Code
	for Scanner.Scan() {
		TotalOutput = append(TotalOutput, Scanner.Text())
	}

	//Memory address needs to increment by 4 every time
	PC_Address = 96
	//BREAK conditional
	if_break = false

	//This for loop grabs the total amount of lines in input and outputs the according amount of dissemabled LEGv8 Code -Hunter
	for _, line := range TotalOutput {
		binaryString = line[0:32]

		//Checks if binaryString is greater then 11
		if len(binaryString) > 11 {

			//If Break operand has been called
			if if_break == true {

				//If sign value is equal to 1 reverse two the string otherwise its positive
				if binaryString[0:1] == "1" {
					writeToOutput = fmt.Sprintf("%s\t\t\t%d\t\t-%d", binaryString, PC_Address, binaryStringToInt(reverseTwosComplement(binaryString[1:32])))
					//Stores data after break into array for later use
					afterBreakData = append(afterBreakData, -1*int32(binaryStringToInt(reverseTwosComplement(binaryString[1:32]))))
				} else {
					writeToOutput = fmt.Sprintf("%s\t\t\t%d\t\t%d", binaryString, PC_Address, binaryStringToInt(binaryString))
					//Stores data after break into array for later use
					afterBreakData = append(afterBreakData, int32(binaryStringToInt(binaryString)))
				}
			} else {
				//Slices first 11 characters
				opcodeString = binaryString[:11]

				//Converts binary string to decimal
				opcodeValue := binaryStringToInt(opcodeString)

				//call findOpcode and store results in writeToOutput
				writeToOutput = findOpcode(opcodeValue, binaryString, uint32(PC_Address))
			}

			Storage.WriteString(writeToOutput + "\n")
			PC_Address = PC_Address + 4

		} else {
			fmt.Println("ERROR: Input String not correct length")
		}
	}

	//Project part 2 starts here
	//for project 2 says "you will use a WHILE LOOP. While not equal to BREAK"; except Golang doesnt have a WHILE loop so we need a for loop instead
	//We will need a size 32 int array for registers and a size 8 int array for data
	var simString string
	//if_STUR_LDUR := false
	//dataAddress exists just past the final PC address in the simulator file
	//             96 + (total length of array - (total length of array - lines after break)) * 4

	cyclePos := 1
	arrIndex = 0

	for i := 0; i < len(RegArray); i++ {
		RegArray[i] = 0
	}

	//the FOR loop will act as a WHILE loop until the BREAK instruction is reached -Alan
	for {
		//similar to Project 1, take the first string from input and convert first 11 bits to an int
		//except this time we use TotalOutput because it has every 32 bit binary string stored
		simString = TotalOutput[arrIndex]
		opcodeString = simString[:11]
		opcodeValue := binaryStringToInt(opcodeString)
		PC_Address = 96

		//finds opcode again using what we already built from project 1
		writeToOutput = findOpcode(opcodeValue, simString, uint32(arrIndex))

		//it appears that the simulator does not display NOP and UNKNOWN instruction
		if writeToOutput == "SKIP" {
			arrIndex += 1
		} else {
			//Writes cycle line as well as pc address and operand info from opcode func
			writer := fmt.Sprintf("cycle:%d\t%s", cyclePos, writeToOutput)
			Simulator.WriteString(writer + "\n\n" + "registers:\n")

			//displays register data[0-7]
			Simulator.WriteString("r00:\t")
			for i := 0; i < 8; i++ {
				writer = fmt.Sprintf("%d\t", RegArray[i])
				Simulator.WriteString(writer)
			}
			Simulator.WriteString("\n")
			//displays register data[8-15]
			Simulator.WriteString("r08:\t")

			for i := 8; i < 16; i++ {
				writer = fmt.Sprintf("%d\t", RegArray[i])
				Simulator.WriteString(writer)
			}
			Simulator.WriteString("\n")
			//displays register data[16-23]
			Simulator.WriteString("r16:\t")
			for i := 16; i < 24; i++ {
				writer = fmt.Sprintf("%d\t", RegArray[i])
				Simulator.WriteString(writer)
			}
			Simulator.WriteString("\n")
			//displays register data[24-31]
			Simulator.WriteString("r24:\t")
			for i := 24; i < 32; i++ {
				writer = fmt.Sprintf("%d\t", RegArray[i])
				Simulator.WriteString(writer)
			}
			Simulator.WriteString("\n\n")
			Simulator.WriteString("data:\n")

			//if STUR or LDUR is called in program, set bool to true
			//so the next func print out DataArray to simulator file
			//if(opcodeValue == 1984 || opcodeValue == 1986){
			//if_STUR_LDUR = true
			//}
			//if boolean is true, print out DataArray to simulator file

			//STUR will have to change a value within afterBreakData[]
			//LDUR will copy a value from within afterBreakData[]

			// Used for formatting the data into the _sim.txt output file
			iterations := int(math.Ceil(float64(len(afterBreakData)) / 8.0))
			for i := 0; i < iterations; i++ {
				writer = fmt.Sprintf("%d:", DataStartIndex+(i*32))
				Simulator.WriteString(writer)

				for j := 0; j < 8; j++ {

					if len(afterBreakData) > (i*8)+j {
						writer = fmt.Sprintf("\t%d", afterBreakData[(i*8)+j])
						Simulator.WriteString(writer)
					} else {
						writer = fmt.Sprintf("\t0")
						Simulator.WriteString(writer)
					}
				}
				Simulator.WriteString("\n")
			}

			//"WHILE" loop break condition
			if opcodeValue == 2038 {
				break
			}

			Simulator.WriteString("====================\n")

			//update variables
			/* IMPORTANT: The PC address does not get updated since it is calculated
			 *  by adding the arrIndex value * 4 in opcodes.go
			 */
			cyclePos += 1
			arrIndex += 1
		}
	}
}

// Takes a binary string input and returns that binary number in decimal. - Hunter
func binaryStringToInt(tempStr string) uint32 {
	decimalValue, err := strconv.ParseInt(tempStr, 2, 64)

	if err != nil {
		fmt.Println("ERROR with binary string to interger conversion")
	}
	return uint32(decimalValue)
}

// Switch statement that grabs a Opcode Interger value and determines which Opcode is being called. - Hunter

//findOpcode now returns a string that will be written to Output file - Alan

//(update)- (9/6) - writeToOutput now reads CBZ and CBNZ values by adding a couple of more case statements in the opcodes.go... both the imtest_bin.txt and cbtest1.txt files have correct opcodes.- Matt.

func findOpcode(OpcodeInt uint32, binaryString string, memoryAddress uint32) string {
	writeToOutput := ""

	switch {
	// B Op
	case OpcodeInt >= 160 && OpcodeInt <= 191:
		writeToOutput = B_op(binaryString, memoryAddress)
		break

		// AND Op
	case OpcodeInt == 1104:
		writeToOutput = R_Format_op_No_Shamt(binaryString, memoryAddress, 2)
		break

		// ADD Op
	case OpcodeInt == 1112:
		writeToOutput = R_Format_op_No_Shamt(binaryString, memoryAddress, 0)
		break

		// ADDI Op
	case OpcodeInt == 1160 || OpcodeInt == 1161:
		writeToOutput = Immediate_op(binaryString, memoryAddress, 0)
		break

		// ORR Op
	case OpcodeInt == 1360:
		writeToOutput = R_Format_op_No_Shamt(binaryString, memoryAddress, 1)
		break

		// CBZ Op
	case OpcodeInt >= 1440 && OpcodeInt <= 1447:
		writeToOutput = CBZ_op(binaryString, memoryAddress)
		break

		// CBNZ Op
	case OpcodeInt >= 1448 && OpcodeInt <= 1455:
		writeToOutput = CBNZ_op(binaryString, memoryAddress)
		break

		// SUB Op
	case OpcodeInt == 1624:
		writeToOutput = R_Format_op_No_Shamt(binaryString, memoryAddress, 3)
		break

		// SUBI Op
	case OpcodeInt == 1672 || OpcodeInt == 1673:
		writeToOutput = Immediate_op(binaryString, memoryAddress, 1)
		break

		//MOVZ Op
	case OpcodeInt >= 1684 && OpcodeInt <= 1687:
		writeToOutput = MOVZ_op(binaryString, memoryAddress)
		break

		// MOVK Op
	case OpcodeInt >= 1940 && OpcodeInt <= 1943:
		writeToOutput = MOVK_op(binaryString, memoryAddress)
		break

		// LSL Op (Type 0)
	case OpcodeInt == 1691:
		writeToOutput = R_Format_op_Shamt(binaryString, memoryAddress, 0)
		break

		// LSR Op (type 1)
	case OpcodeInt == 1690:
		writeToOutput = R_Format_op_Shamt(binaryString, memoryAddress, 1)
		break

		// ASR Op (type 2)
	case OpcodeInt == 1692:
		writeToOutput = R_Format_op_Shamt(binaryString, memoryAddress, 2)
		break

		// STUR Op
	case OpcodeInt == 1984:
		writeToOutput = LastOps(binaryString, memoryAddress, 0)
		break

		// LDUR Op
	case OpcodeInt == 1986:
		writeToOutput = LastOps(binaryString, memoryAddress, 1)
		break

		// EOR Op
	case OpcodeInt == 1872:
		writeToOutput = R_Format_op_No_Shamt(binaryString, memoryAddress, 4)
		break

		// NOP Op
	case OpcodeInt == 0:
		//Writes NOP format string directly without a function call
		if if_break == false {
			if binaryString[0:32] == "00000000000000000000000000000000" {
				writeToOutput = binaryString[0:32] + "\t\t" + fmt.Sprintf("\t%-3d", memoryAddress) + "\t\tNOP"
			}
		} else {
			writeToOutput = "SKIP"
		}

		break

	// Break Op
	case OpcodeInt == 2038:
		if if_break == false {
			//This outputs the break format directly without a function
			writeToOutput = binaryString[0:1] + " " + binaryString[1:6] + " " + binaryString[6:11] + " " + binaryString[11:16] + " " + binaryString[16:21] + " " + binaryString[21:26] + " " + binaryString[26:32] + "\t\t" + fmt.Sprintf("%-3d", memoryAddress) + "\t\tBREAK"

			//Once break is called no more operands should be found, this goes back to main
			DataStartIndex = int(memoryAddress + 4)
			if_break = true
		} else {
			PC_Address = PC_Address + int(memoryAddress)*4
			writeToOutput = fmt.Sprintf("%d BREAK", PC_Address)
		}

		break

		//No Opp Code Found
	default:
		if if_break == false {
			writeToOutput = binaryString[0:11] + " " + binaryString[11:32] + "\t" + fmt.Sprintf("\t\t%-3d", memoryAddress) + "\t\tUnknown Instruction"
		} else {
			writeToOutput = "SKIP"
		}

		break
	}

	return writeToOutput
}
