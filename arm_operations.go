//This go file was made to seperate all the individual OP functions from the main file. -Hunter
/*3) A checker needs to have every available op for every instruction type to compare op codes to. Once a matching code has been found and instruction type has been discovered, send 32 bit string to an decoder that will translate binary to the ARMv8 per instruction type. -The decoder will have to split up the 32 bit string a specific way based on the instruction set.*/

package main

import (
	//"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var returnString string

// Takes binary string input and outputs corresponing reverse twos complement string. -Hunter
// EX. 110100 --> 001100
func reverseTwosComplement(binString string) string {
	//Converts to binary
	//Grabs length of binary before and after due to when converting to decimal it cuts off the 0s. For example 0000000101 would turn  into 101. This length is used later in a for loop to fix this.
	length_before := len(binString)
	decimal, err := strconv.ParseInt(binString, 2, 64)
	if err != nil {
		fmt.Println("ERROR: Could not perform reverseTwosComplement")
		return binString
	}

	// Subtract 1
	decimal -= 1

	// Convert back to binary
	binString = strconv.FormatInt(decimal, 2)

	//Length of string after conversion
	length_after := len(binString)

	//This adds the 0s back by how much length was lost
	for i := 1; i <= length_before-length_after; i++ {
		binString = "0" + binString
	}

	//Flips all bits
	var flip strings.Builder
	for _, char := range binString {
		if char == '0' {
			flip.WriteRune('1')
		} else if char == '1' {
			flip.WriteRune('0')
		} else {
			flip.WriteRune(char)
		}
	}
	return flip.String()
}

// MOVZ Operand -Hunter
// func returns string to findOpcode -Alan
func MOVZ_op(binaryString string, memoryAddress uint32) string {
	//First 9 are OP code, next 2 are shift code (0, 16, 32, 48), next 16 are immidiate value, last 5 are Rd

	//Seperates and creates strings of seperation
	shift_code := binaryString[9:11]
	immidiate := binaryString[11:27]
	Rd := binaryString[27:32]

	//Converts immidiate to decimal
	decimal_immidiate, err := strconv.ParseInt(immidiate, 2, 64)

	//Converts Rd to decimal
	decimal_Rd, err := strconv.ParseInt(Rd, 2, 64)

	if err != nil {
		fmt.Println("ERROR: Could not perform MOVZ_op")
		os.Exit(-1)
	}

	var decimal_shift int32
	//Determines shift value from shift code
	switch {
	case shift_code == "00":
		decimal_shift = 0
	case shift_code == "01":
		decimal_shift = 16
	case shift_code == "10":
		decimal_shift = 32
	case shift_code == "11":
		decimal_shift = 48
	}

	//returnString stores the string that will be written into Output file. fmt.Sprintf allows for the formatting of a string using multiple value types with specific formatting option (%s = String, %d = decimal int)
	if if_break == false {
		returnString = fmt.Sprintf("%s %s %s %s\t\t\t%-3d\t\t%s%d, %d, %s, %d", binaryString[0:9], binaryString[9:11], binaryString[11:27], binaryString[27:32], memoryAddress, "MOVZ\tR", decimal_Rd, decimal_immidiate, "LSL", decimal_shift)
	} else {

		PC_Address = PC_Address + int(memoryAddress)*4
		returnString = fmt.Sprintf("%d MOVZ", PC_Address)
	}

	return returnString
}

// MOVK Operand -Hunter
func MOVK_op(binaryString string, memoryAddress uint32) string {
	//First 9 are OP code, next 2 are shift code (0, 16, 32, 48), next 16 are immidiate value, last 5 are Rd

	//Seperates and creates strings of seperation
	shift_code := binaryString[9:11]
	immidiate := binaryString[11:27]
	Rd := binaryString[27:32]

	//Converts immidiate to decimal
	decimal_immidiate, err := strconv.ParseInt(immidiate, 2, 64)

	//Converts Rd to decimal
	decimal_Rd, err := strconv.ParseInt(Rd, 2, 64)

	if err != nil {
		fmt.Println("ERROR: Could not perform MOVK_op")
		os.Exit(-1)
	}

	var decimal_shift int32
	//Determines shift value from shift code
	switch {
	case shift_code == "00":
		decimal_shift = 0
	case shift_code == "01":
		decimal_shift = 16
	case shift_code == "10":
		decimal_shift = 32
	case shift_code == "11":
		decimal_shift = 48
	}

	//returnString stores the string that will be written into Output file. fmt.Sprintf allows for the formatting of a string using multiple value types with specific formatting option (%s = String, %d = decimal int)
	if if_break == false {
		returnString = fmt.Sprintf("%s %s %s %s\t\t\t%-3d\t\t%s%d, %d, %s, %d", binaryString[0:9], binaryString[9:11], binaryString[11:27], binaryString[27:32], memoryAddress, "MOVK\tR", decimal_Rd, decimal_immidiate, "LSL", decimal_shift)
	} else {

		PC_Address = PC_Address + int(memoryAddress)*4
		returnString = fmt.Sprintf("MOVK %d", PC_Address)
	}

	return returnString
}

// Branch Operand -Hunter
func B_op(binaryString string, memoryAddress uint32) string {
	//Grabs only the sign bit to determine if negative or positve
	sign_bit := binaryString[6:7]

	var reverse_offset string
	var decimal_offset int64
	var err error
	//Grabs the offset value without sign bit to calculate branch amount
	tempoffset := binaryString[7:32]

	//Checks if sign bit is 1 and reverse twos compliment everything past the sign bit.
	if sign_bit == "1" {
		//Reverse twos compliments offset
		reverse_offset = reverseTwosComplement(tempoffset)
		//Converts to decimal
		decimal_offset, err = strconv.ParseInt(reverse_offset, 2, 64)

		if err != nil {
			fmt.Println("ERROR: Could not perform B_op")
			os.Exit(-1)
		}
		//Multiples by -1 to make value negative
		decimal_offset = decimal_offset * -1
	} else {
		//Converts to decimal
		decimal_offset, err = strconv.ParseInt(tempoffset, 2, 64)
	}

	if if_break == false {
		returnString = fmt.Sprintf("%s %s\t\t\t%-3d\t\t%s%d", binaryString[0:6], binaryString[6:32], memoryAddress, "B\t\t\t#", decimal_offset)
	} else {

		PC_Address = PC_Address + int(memoryAddress)*4
		arrIndex = arrIndex + int(decimal_offset) - 1
		returnString = fmt.Sprintf("%d\tB\t#%d", PC_Address, decimal_offset)
	}

	return returnString
}

// This function works with ADD ORR EOR AND SUB, since there format is identical all that changes is the outputed operand name -Hunter
func R_Format_op_No_Shamt(binaryString string, memoryAddress uint32, RFormatType uint32) string {
	//Finding sections of binary and turning to decimal
	decimal_R1, err := strconv.ParseInt(binaryString[22:27], 2, 64)
	decimal_R3, err := strconv.ParseInt(binaryString[27:32], 2, 64)
	decimal_R2, err := strconv.ParseInt(binaryString[11:16], 2, 64)

	var RFormatType_String string
	//Finds what opperand type it is from the RFormatType input
	switch RFormatType {
	case 0:
		RFormatType_String = "ADD"

		if if_break == true {
			RegArray[decimal_R3] = RegArray[decimal_R2] + RegArray[decimal_R1]
		}
		break

	case 1:
		RFormatType_String = "ORR"
		if if_break == true {
			RegArray[decimal_R3] = RegArray[decimal_R2] | RegArray[decimal_R1]
		}
		break

	case 2:
		RFormatType_String = "AND"
		if if_break == true {
			RegArray[decimal_R3] = RegArray[decimal_R2] & RegArray[decimal_R1]
		}
		break

	case 3:
		RFormatType_String = "SUB"
		if if_break == true {
			RegArray[decimal_R3] = RegArray[decimal_R1] - RegArray[decimal_R2]
		}
		break

	case 4:
		RFormatType_String = "EOR"
		if if_break == true {
			RegArray[decimal_R3] = RegArray[decimal_R2] ^ RegArray[decimal_R1]
		}
		break

	default:
		fmt.Println("ERROR: Could not find R_Format Type")
		os.Exit(-1)
		break
	}
	if err != nil {
		fmt.Println("ERROR: Could not perform R_Format_op_No_Shamt")
		os.Exit(-1)
	}

	if if_break == false {
		returnString = fmt.Sprintf("%s %s %s %s %s\t\t%-3d\t\t%s\t\tR%d, R%d, R%d", binaryString[0:11], binaryString[11:16], binaryString[16:22], binaryString[22:27], binaryString[27:32], memoryAddress, RFormatType_String, decimal_R3, decimal_R1, decimal_R2)
	} else {

		PC_Address = PC_Address + int(memoryAddress)*4
		returnString = fmt.Sprintf("%d\t%s\tR%d, R%d, R%d", PC_Address, RFormatType_String, decimal_R3, decimal_R1, decimal_R2)
	}

	return returnString
}

// CBZ op func (it doesn't have the same format as the R type)
// CBZ op takes the 2's complement of the offset to determine if its positive or negative -Alan
func CBZ_op(binaryString string, memoryAddress uint32) string {
	//Grabs only the sign bit to determine if negative or positve
	sign_bit := binaryString[8:9]

	var reverse_offset string
	var decimal_offset int64
	var err error
	//Grabs the offset value without sign bit to calculate branch amount
	tempoffset := binaryString[9:27]

	//Checks if sign bit is 1 and reverse twos compliment everything past the sign bit.
	if sign_bit == "1" {
		//Reverse twos compliments offset
		reverse_offset = reverseTwosComplement(tempoffset)
		//Converts to decimal
		decimal_offset, err = strconv.ParseInt(reverse_offset, 2, 64)

		if err != nil {
			fmt.Println("ERROR: Could not perform CBZ_op")
			os.Exit(-1)
		}
		//Multiples by -1 to make value negative
		decimal_offset = decimal_offset * -1
	} else {
		//Converts to decimal
		decimal_offset, err = strconv.ParseInt(tempoffset, 2, 64)
	}

	//Converts the last 5 bits into a decimal number
	decimal_Conditional, err := strconv.ParseInt(binaryString[27:32], 2, 64)

	if if_break == false {
		returnString = fmt.Sprintf("%s %s %s\t\t\t%-3d\t\t%s\t\t%c%d%s%d", binaryString[0:8], binaryString[8:27], binaryString[27:32], memoryAddress, "CBZ", 'R', decimal_Conditional, ", #", decimal_offset)
	} else {

		PC_Address = PC_Address + int(memoryAddress)*4
		if RegArray[decimal_Conditional] == 0 {
			arrIndex = arrIndex + int(decimal_offset) - 1
		}

		returnString = fmt.Sprintf("%d\tCBZ\tR%d, #%d", PC_Address, decimal_Conditional, decimal_offset)
	}

	return returnString
}

// CBNZ is essentially the same as CBZ, both use the 2's complement to output a signed number -Alan
func CBNZ_op(binaryString string, memoryAddress uint32) string {
	//Grabs only the sign bit to determine if negative or positve
	sign_bit := binaryString[8:9]

	var reverse_offset string
	var decimal_offset int64
	var err error
	//Grabs the offset value without sign bit to calculate branch amount
	tempoffset := binaryString[9:27]

	//Checks if sign bit is 1 and reverse twos compliment everything past the sign bit.
	if sign_bit == "1" {
		//Reverse twos compliments offset
		reverse_offset = reverseTwosComplement(tempoffset)
		//Converts to decimal
		decimal_offset, err = strconv.ParseInt(reverse_offset, 2, 64)

		if err != nil {
			fmt.Println("ERROR: Could not perform CBNZ_op")
			os.Exit(-1)
		}
		//Multiples by -1 to make value negative
		decimal_offset = decimal_offset * -1
	} else {
		//Converts to decimal
		decimal_offset, err = strconv.ParseInt(tempoffset, 2, 64)
	}

	//Converts the last 5 bits into a decimal number
	decimal_Conditional, err := strconv.ParseInt(binaryString[27:32], 2, 64)

	if if_break == false {
		returnString = fmt.Sprintf("%s %s %s\t\t\t%-3d\t\t%s\t%c%d%s%d", binaryString[0:8], binaryString[8:27], binaryString[27:32], memoryAddress, "CBNZ", 'R', decimal_Conditional, ", #", decimal_offset)
	} else {

		if RegArray[decimal_Conditional] != 0 {
			arrIndex = arrIndex + int(decimal_offset) - 1
		}

		PC_Address = PC_Address + int(memoryAddress)*4
		returnString = fmt.Sprintf("%d\tCBNZ\tR%d, #%d", PC_Address, decimal_Conditional, decimal_offset)
	}

	return returnString
}

// Outputs the last three R type instructions -Alan
func R_Format_op_Shamt(binaryString string, memoryAddress uint32, RFormatType uint32) string {
	//Finding sections of binary and turning to decimal
	//LSL, LSR, and ASR only uses R3(Rd) & R1(Rn)
	decimal_Rn, err := strconv.ParseInt(binaryString[22:27], 2, 64)
	decimal_Rd, err := strconv.ParseInt(binaryString[27:32], 2, 64)
	decimal_Shamt, err := strconv.ParseInt(binaryString[16:22], 2, 64)

	var RFormatType_String string
	//Finds what opperand type it is from the RFormatType input
	switch RFormatType {
	case 0:
		RFormatType_String = "LSL"
		if if_break == true {
			RegArray[decimal_Rd] = RegArray[decimal_Rn] << int(decimal_Shamt)
		}
		break

	case 1:
		RFormatType_String = "LSR"
		if if_break == true {
			if RegArray[decimal_Rn] < 0 {
				RegArray[decimal_Rn] = int(binaryStringToInt(reverseTwosComplement(fmt.Sprintf("%032b", int(math.Abs(float64(RegArray[decimal_Rd])))))))
				RegArray[decimal_Rd] = RegArray[decimal_Rn] >> int(decimal_Shamt)
			} else {
				RegArray[decimal_Rd] = RegArray[decimal_Rn] >> int(decimal_Shamt)
			}
		}
		break

	case 2:
		RFormatType_String = "ASR"
		if if_break == true {
			if RegArray[decimal_Rn] < 0 {
				RegArray[decimal_Rd] = int(math.Abs(float64(RegArray[decimal_Rn]))) >> int(decimal_Shamt)
				RegArray[decimal_Rd] = RegArray[decimal_Rd] * -1
			} else {
				RegArray[decimal_Rd] = RegArray[decimal_Rn] >> int(decimal_Shamt)
			}
		}
		break

	default:
		fmt.Println("ERROR: Could not find R_Format Type")
		os.Exit(-1)
		break
	}
	if err != nil {
		fmt.Println("ERROR: Could not perform R_Format_op_Shamt")
		os.Exit(-1)
	}

	if if_break == false {
		returnString = fmt.Sprintf("%s %s %s %s %s\t\t%-3d\t\t%s\t\tR%d, R%d, #%d", binaryString[0:11], binaryString[11:16], binaryString[16:22], binaryString[22:27], binaryString[27:32], memoryAddress, RFormatType_String, decimal_Rd, decimal_Rn, decimal_Shamt)
	} else {

		PC_Address = PC_Address + int(memoryAddress)*4
		returnString = fmt.Sprintf("%d\t%s\tR%d, R%d, #%d", PC_Address, RFormatType_String, decimal_Rd, decimal_Rn, decimal_Shamt)
	}

	return returnString
}

// finds the opcodes for STUR and LDUR(10/9/23).
func LastOps(binaryString string, memoryAddress uint32, DFormatType uint32) string {
	RnDecimal, err := strconv.ParseInt(binaryString[22:27], 2, 64) //converts to decimal.
	RDDecimal, err := strconv.ParseInt(binaryString[27:32], 2, 64) //converts to decimal.
	Address, err := strconv.ParseInt(binaryString[11:20], 2, 64)   //converts to decimal.

	if err != nil { //if err != nil, print error and exit.
		fmt.Println("Error: conversion error")
		os.Exit(-1)
	}
	var other_String string //establishes a string variable.
	switch DFormatType {    //establishes a switch to grab D-Format and determine if the opcode is STUR or LDUR.

	case 0:
		other_String = "STUR"
		if if_break == true {
			//Checks if there is avaiable data, if there isn't it increases size, if there is it stores normally, if there is an error it prints an error.
			if len(afterBreakData) < (RegArray[RnDecimal]+int(Address)*4-(DataStartIndex))/4 {
				var size_Difference = ((RegArray[RnDecimal] + int(Address)*4 - (DataStartIndex)) / 4) - len(afterBreakData)
				for i := 0; i < size_Difference+1; i++ {
					afterBreakData = append(afterBreakData, 0)
				}
				afterBreakData[(RegArray[RnDecimal]+int(Address)*4-(DataStartIndex))/4] = int32(RegArray[RDDecimal])
			} else if 0 > (RegArray[RnDecimal]+int(Address)*4-(DataStartIndex))/4 {
				fmt.Println("ERROR: Check STUR data address, exceeds Data Memory size")
				os.Exit(-300)
			} else {
				afterBreakData[(RegArray[RnDecimal]+int(Address)*4-(DataStartIndex))/4] = int32(RegArray[RDDecimal])
			}

		}
		break

	case 1:
		other_String = "LDUR"
		if if_break == true {
			// fmt.Println(DataStartIndex);
			//afterBreakData[(RegArray[RnDecimal]+int(Address)*4-(DataStartIndex-4))/4] = int32(RegArray[RDDecimal]);

			//Checks if there is avaiable data, if there isn't it increases size, if there is it stores normally, if there is an error it prints an error.
			if len(afterBreakData) < (RegArray[RnDecimal]+int(Address)*4-(DataStartIndex))/4 {
				var size_Difference = ((RegArray[RnDecimal] + int(Address)*4 - (DataStartIndex)) / 4) - len(afterBreakData)
				for i := 0; i < size_Difference+1; i++ {
					afterBreakData = append(afterBreakData, 0)
				}
				RegArray[RDDecimal] = 0
			} else if 0 > (RegArray[RnDecimal]+int(Address)*4-(DataStartIndex))/4 {
				fmt.Println("ERROR: Check LDUR data address, exceeds Data Memory size")
				os.Exit(-300)
			} else {
				RegArray[RDDecimal] = int(afterBreakData[(RegArray[RnDecimal]+int(Address)*4-(DataStartIndex))/4])
			}
		}
		break

	default:
		fmt.Println("ERROR couldnt find correct OP")
		os.Exit(-1)
		break
	}

	D, err := strconv.ParseInt(binaryString[11:20], 2, 64)
	//var tempString string
	if if_break == false {
		returnString = fmt.Sprintf("%s %s %s %s %s\t\t%-3d\t\t%s\tR%d, [R%d, #%d]",
			binaryString[0:11], binaryString[11:20], binaryString[20:22], binaryString[22:27], binaryString[27:32], memoryAddress, other_String, RDDecimal, RnDecimal, D)
	} else {

		PC_Address = PC_Address + int(memoryAddress)*4
		returnString = fmt.Sprintf("%d\t%s\tR%d, [R%d, #%d]", PC_Address, other_String, RDDecimal, RnDecimal, D)
	}

	return returnString
}

func Immediate_op(binaryString string, memoryAddress uint32, IFormatType uint32) string {

	decimal_Rn, err := strconv.ParseInt(binaryString[22:27], 2, 64)
	decimal_Rd, err := strconv.ParseInt(binaryString[27:32], 2, 64)
	decimal_offset, err := strconv.ParseInt(binaryString[10:22], 2, 64)

	var IFormatType_String string
	//Finds what opperand type it is from the RFormatType input
	switch IFormatType {
	case 0:
		IFormatType_String = "ADDI"
		if if_break == true {
			RegArray[decimal_Rd] = RegArray[decimal_Rn] + int(decimal_offset)
		}
		break

	case 1:
		IFormatType_String = "SUBI"
		if if_break == true {
			RegArray[decimal_Rd] = RegArray[decimal_Rn] - int(decimal_offset)
		}
		break

	default:
		fmt.Println("ERROR: Could not find I_Format Type")
		os.Exit(-1)
		break
	}
	if err != nil {
		fmt.Println("ERROR: Could not perform Immidiate_op")
		os.Exit(-1)
	}

	if if_break == false {
		returnString = fmt.Sprintf("%s %s %s %s\t\t\t%-3d\t\t%s\tR%d, R%d, #%d", binaryString[0:10], binaryString[10:22], binaryString[22:27], binaryString[27:32], memoryAddress, IFormatType_String, decimal_Rd, decimal_Rn, decimal_offset)
	} else {

		PC_Address = PC_Address + int(memoryAddress)*4
		returnString = fmt.Sprintf("%d\t%s\tR%d, R%d, #%d", PC_Address, IFormatType_String, decimal_Rd, decimal_Rn, decimal_offset)
	}

	return returnString
}
