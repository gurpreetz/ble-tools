package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/currantlabs/gatt"
)

// XMLCharProperties represents the BLE characteristic properties from the xml file
type XMLCharProperties struct {
	Broadcast            string
	Read                 string
	WriteWithoutResponse string
	Write                string
	Notify               string
	Indicate             string
	SignedWrite          string
	Extended             string
	bitMask              gatt.Property
}

// XMLCharacteristic represents the BLE characteristic information from the xml file
type XMLCharacteristic struct {
	CharName    string `xml:"name,attr"`
	CharID      string `xml:"uuid,attr"`
	Requirement string
	Properties  XMLCharProperties
}

// XMLService represents the BLE service information from the xml file
type XMLService struct {
	ServiceName string              `xml:"name,attr"`
	ServiceID   string              `xml:"uuid,attr"`
	CharList    []XMLCharacteristic `xml:"characteristic"`
	numChars    int
}

// XMLDevice represents the BLE Device information from the xml file
type XMLDevice struct {
	XMLName     xml.Name     `xml:"device"`
	DeviceName  string       `xml:"name,attr"`
	ServiceList []XMLService `xml:"service"`
	numServices int
}

const mandatory = "Mandatory"
const excluded = "Excluded"

// getProperties gets a bitmap of the characteristic properties
func getProperties(char *XMLCharacteristic) {
	fmt.Print("\t    ")
	if strings.Compare(char.Properties.Broadcast, mandatory) == 0 {
		fmt.Print("broadcast ")
		char.Properties.bitMask |= gatt.CharBroadcast
	}
	if strings.Compare(char.Properties.Read, mandatory) == 0 {
		fmt.Print("read ")
		char.Properties.bitMask |= gatt.CharRead
	}
	if strings.Compare(char.Properties.WriteWithoutResponse, mandatory) == 0 {
		fmt.Print("writeWithoutResponse ")
		char.Properties.bitMask |= gatt.CharWriteNR
	}
	if strings.Compare(char.Properties.Write, mandatory) == 0 {
		fmt.Print("write ")
		char.Properties.bitMask |= gatt.CharWrite
	}
	if strings.Compare(char.Properties.Notify, mandatory) == 0 {
		fmt.Print("notify ")
		char.Properties.bitMask |= gatt.CharNotify
	}
	if strings.Compare(char.Properties.Indicate, mandatory) == 0 {
		fmt.Print("indicate ")
		char.Properties.bitMask |= gatt.CharIndicate
	}
	if strings.Compare(char.Properties.SignedWrite, mandatory) == 0 {
		fmt.Print("signedWrite ")
		char.Properties.bitMask |= gatt.CharSignedWrite
	}
	if strings.Compare(char.Properties.Extended, mandatory) == 0 {
		fmt.Print("extended  ")
		char.Properties.bitMask |= gatt.CharExtended
	}
	fmt.Println()
}

// xmlSetProperties sets the xml characteristic properties based on the bitmap
func xmlSetProperties(prop gatt.Property) *XMLCharProperties {
	var xmlProp XMLCharProperties

	if (prop & gatt.CharBroadcast) != 0 {
		xmlProp.Broadcast = mandatory
	} else {
		xmlProp.Broadcast = excluded
	}
	if (prop & gatt.CharRead) != 0 {
		xmlProp.Read = mandatory
	} else {
		xmlProp.Read = excluded
	}
	if (prop & gatt.CharWriteNR) != 0 {
		xmlProp.WriteWithoutResponse = mandatory
	} else {
		xmlProp.WriteWithoutResponse = excluded
	}
	if (prop & gatt.CharWrite) != 0 {
		xmlProp.Write = mandatory
	} else {
		xmlProp.Write = excluded
	}
	if (prop & gatt.CharNotify) != 0 {
		xmlProp.Notify = mandatory
	} else {
		xmlProp.Notify = excluded
	}
	if (prop & gatt.CharIndicate) != 0 {
		xmlProp.Indicate = mandatory
	} else {
		xmlProp.Indicate = excluded
	}
	if (prop & gatt.CharSignedWrite) != 0 {
		xmlProp.SignedWrite = mandatory
	} else {
		xmlProp.SignedWrite = excluded
	}
	if (prop & gatt.CharExtended) != 0 {
		xmlProp.Extended = mandatory
	} else {
		xmlProp.Extended = excluded
	}
	return &xmlProp
}

// xmlAppendCharInfo appends the current characteristic information to the xml being generated
func xmlAppendCharInfo(charName string, charUUID string, prop gatt.Property) *XMLCharacteristic {
	var xmlChar XMLCharacteristic

	xmlProp := xmlSetProperties(prop)

	xmlChar.CharName = charName
	xmlChar.CharID = charUUID
	xmlChar.Requirement = "mandatory"
	xmlChar.Properties = *xmlProp

	return &xmlChar

}

// xmlAppendSvcInfo appends the current service information to the xml being generated
func xmlAppendSvcInfo(device *XMLDevice, svcName string, svcUUID string, charList []XMLCharacteristic) *XMLService {
	var xmlSvc XMLService

	xmlSvc.ServiceName = svcName
	xmlSvc.ServiceID = svcUUID
	xmlSvc.CharList = charList

	return &xmlSvc
}

// xmlFindService searches for a service, by UUID, in a given xml parsed device
func xmlFindService(device *XMLDevice, svcID string) (bool, *XMLService) {
	for _, s := range device.ServiceList {
		if svcID == s.ServiceID {
			return true, &s
		}
	}
	return false, nil
}

// xmlFindChar searches for a characteristic, by UUID, in a given xml parsed service
func xmlFindChar(svc *XMLService, charID string) (bool, *XMLCharacteristic) {
	for _, c := range svc.CharList {
		if charID == c.CharID {
			return true, &c
		}
	}
	return false, nil
}

// xmlShowDeviceSummary displayes the summary of a device parsed from an xml file
func xmlShowDeviceSummary(device *XMLDevice) {

	var svcName string
	fmt.Println()
	fmt.Println("***** DEVICE SUMMARY *****")
	fmt.Println(device.DeviceName, "has", device.numServices, "services")

	for _, s := range device.ServiceList {
		if len(s.ServiceName) != 0 {
			svcName = s.ServiceName
		} else {
			svcName = s.ServiceID
		}
		fmt.Println("  ", svcName, "has", s.numChars, "characteristic(s)")
	}
	fmt.Println("**************************")
}

// xmlOutDeviceInfo creates an xml file based on the BLE device the tool is connected to
func xmlOutDeviceInfo(dev *XMLDevice) {
	var dirName = "XmlOutputs"

	_, err := os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			if os.Mkdir(dirName, 0777) != nil {
				panic("Unable to create directory" + dirName)
			}
		}
	}

	s := []string{dirName, "/", dev.DeviceName, ".xml"}
	var xmlFile = strings.Join(s, "")

	fmt.Println("XML Output created in file", xmlFile)

	f, err := os.Create(xmlFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	output, err := xml.MarshalIndent(dev, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	f.Write([]byte(xml.Header))
	f.Write(output)

	fmt.Println()
}

// xmlGetServices parsed an xml file to create a representation of the device in memory
func xmlGetServices(fileName string) *XMLDevice {
	var device XMLDevice

	if len(fileName) == 0 {
		fmt.Println(" Please specify a file to open")
		os.Exit(0)
	}
	xmlFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(b, &device)

	fmt.Println("\nReading Device File for Device ", device.DeviceName)

	for svcIdx, s := range device.ServiceList {
		fmt.Println("Service: ", s.ServiceID, "(", s.ServiceName, ")")
		device.numServices++
		for idx, c := range s.CharList {
			fmt.Println("\tCharacteristic", c.CharID, "(", c.CharName, ")")
			device.ServiceList[svcIdx].numChars++
			getProperties(&s.CharList[idx])
		}
	}
	xmlShowDeviceSummary(&device)

	return &device
}
