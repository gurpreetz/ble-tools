package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Songmu/prompter"
)

func main() {
	scanCommand := flag.NewFlagSet("scan", flag.ExitOnError)
	scanTimeoutFlag := scanCommand.Duration("timeout", 12*time.Second, "scan `timeout` duration in seconds")

	connectCommand := flag.NewFlagSet("connect", flag.ExitOnError)
	connectDeviceFlag := connectCommand.String("device", "", "BLE `Device Name`")
	connectIDFlag := connectCommand.String("id", "", "Last 3 hex bytes of `mfg data` to uniquely identify device")
	connectXMLOutFlag := connectCommand.Bool("xmlOut", false, "generate an xml output")

	readFileCommand := flag.NewFlagSet("read", flag.ExitOnError)
	readXMLFileFlag := readFileCommand.String("file", "", "`xml file` to be parsed")

	compareFileCommand := flag.NewFlagSet("compare", flag.ExitOnError)
	compareDeviceFlag := compareFileCommand.String("device", "", "BLE `Device Name`")
	compareIDFlag := compareFileCommand.String("id", "", "Last 3 hex bytes of `mfg data` to uniquely identify device")
	compareFileFlag := compareFileCommand.String("file", "", "`XML file` to compare against")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [COMMAND] [<options>]\n", os.Args[0])
		fmt.Println("scan")
		scanCommand.PrintDefaults()
		fmt.Println("connect")
		connectCommand.PrintDefaults()
		fmt.Println("read")
		readFileCommand.PrintDefaults()
		fmt.Println("compare")
		compareFileCommand.PrintDefaults()
	}
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "scan":
		scanCommand.Parse(os.Args[2:])

	case "connect":
		connectCommand.Parse(os.Args[2:])

	case "read":
		readFileCommand.Parse(os.Args[2:])

	case "compare":
		compareFileCommand.Parse(os.Args[2:])
	}

	if scanCommand.Parsed() {
		if *scanTimeoutFlag < time.Second {
			fmt.Println("Please enter a scan value of atleast 1s")
			return
		}
		bleScanDevices(*scanTimeoutFlag)
	}

	if readFileCommand.Parsed() {
		if *readXMLFileFlag == "" {
			fmt.Println("Please enter the file to read")
			readFileCommand.PrintDefaults()
			return
		}
		xmlGetServices(*readXMLFileFlag)
	}

	if connectCommand.Parsed() {
		if *connectDeviceFlag == "" {
			fmt.Println("Please enter the name of a device to connect to")
			connectCommand.PrintDefaults()
			return
		}
		fmt.Println("Device :", *connectDeviceFlag, "\tID : ", *connectIDFlag)
		deviceName = *connectDeviceFlag
		if *connectXMLOutFlag == true {
			bleReadDeviceXML(*connectIDFlag, deviceName)
		} else {
			bleReadDevice(*connectIDFlag, deviceName)
		}
	}

	if compareFileCommand.Parsed() {
		if *compareDeviceFlag == "" {
			fmt.Println("Please enter the name of a device to connect to")
			compareFileCommand.PrintDefaults()
			return
		}
		if *compareFileFlag == "" {
			fmt.Println("Please enter the file name to compare with")
			compareFileCommand.PrintDefaults()
			return
		}
		deviceName = *compareDeviceFlag
		bleCompareDevice(*compareIDFlag, deviceName, *compareFileFlag)
	}
}

// cmdGetDeviceConnectId gets the ID of the device to connect to
func cmdGetDeviceConnectID() uint8 {
	var id uint8
	validInput := false
	for validInput == false {
		fmt.Print("Enter Device to connect to: ")
		fmt.Scan(&id)
		if id < scanResultTotal {
			validInput = true
		} else {
			fmt.Println("Please enter a value between 0 and ", (scanResultTotal - 1))
		}
	}
	return id
}

// cmdGetXmlStatus gets the user's input on whether xml should be generated
func cmdGetXMLStatus() bool {
	if prompter.YN("Generate XML after discovery?", false) {
		return true
	}
	return false
}
