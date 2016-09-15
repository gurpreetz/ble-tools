package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/currantlabs/gatt"
	"github.com/currantlabs/gatt/examples/option"
)

var done = make(chan struct{})
var connected = make(chan bool)
var deviceName string
var macID []byte
var isCmpMode = false
var isXMLMode = false
var device *XMLDevice

const maxScanResult uint8 = 30
const maxTimeoutTime time.Duration = 15 * time.Second

var scanResultTotal uint8

// ScanMapResult represents a map entry of a scanned device
type ScanMapResult struct {
	peripheralName string
	peripheral     gatt.Peripheral
	scanResultNum  uint8
}

// ScanListResult represents a list entry of a scanned device
type ScanListResult struct {
	peripheralName string
	peripheral     gatt.Peripheral
}

var scanMap map[string]ScanMapResult
var scanList []ScanListResult

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println()
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

// onPeriphDiscovered Checks the peripheral that is discovered and connects to the correct peripheral
func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if strings.ToUpper(a.LocalName) != strings.ToUpper(deviceName) {
		return
	}
	lenMfgData := len(a.ManufacturerData)

	if len(macID) != 0 && len(a.ManufacturerData) == 0 {
		return
	}

	if (len(a.ManufacturerData) != 0) && (len(macID) != 0) {
		macIDAdv := a.ManufacturerData[lenMfgData-3 : lenMfgData]

		// Compare tail of macIdAdv with tail of macId
		if bytes.Equal(macIDAdv, (macID)) == false {
			return
		}
	}
	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	fmt.Println("  Local Name        =", a.LocalName)
	fmt.Println("  TX Power Level    =", a.TxPowerLevel)
	fmt.Println("  Manufacturer Data =", a.ManufacturerData)
	fmt.Println("  Service Data      =", a.ServiceData)
	fmt.Println("")

	fmt.Println("connecting.... ")
	p.Device().Connect(p)
}

// displayScanResults lists out the results of a passive scan
func displayScanResults() {
	if scanResultTotal == 0 {
		fmt.Println("No Devices discovered")
		return
	}
	fmt.Println("Following Devices discovered:")
	fmt.Println("\tIndex \t Device Name")
	for idx, value := range scanList {
		fmt.Println("\t", idx, "\t", value.peripheralName)
	}

	devID := cmdGetDeviceConnectID()
	deviceName = scanList[devID].peripheralName
	isXMLMode = cmdGetXMLStatus()
	fmt.Println("Connecting to ", devID, "....", scanList[devID].peripheralName)

	p := scanList[devID].peripheral
	p.Device().Connect(p)
}

// onScanPeriphDiscovered Adds discovered peripherals to a map and a list for easy connectivity
func onScanPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {

	if _, ok := scanMap[p.ID()]; ok {
		// already discovered this device; do nothing
		return
	}
	var devName string

	if len(a.LocalName) != 0 {
		devName = a.LocalName
	} else {
		devName = "Unknown"
	}
	if len(a.ManufacturerData) != 0 {
		macIDAdv := string(a.ManufacturerData[len(a.ManufacturerData)-3:])
		hexSt := fmt.Sprintf("%x", macIDAdv)
		devName = devName + "-" + hexSt
	} else {
		uuid := p.ID()
		uuidTail := uuid[len(uuid)-6:]
		devName = devName + "-" + uuidTail
	}

	scanMap[p.ID()] = ScanMapResult{peripheralName: devName, scanResultNum: scanResultTotal, peripheral: p}
	scanResultTotal++

	foundDev := ScanListResult{peripheralName: devName, peripheral: p}

	scanList = append(scanList, foundDev)

	if scanResultTotal == maxScanResult {
		p.Device().StopScanning()
		displayScanResults()
		return
	}
}

func testCharWrite(p gatt.Peripheral, ch *gatt.Characteristic) {
	fmt.Println("testing char " + ch.UUID().String())
	if (ch.Properties() & gatt.CharWrite) == 0 {
		fmt.Println("Char is not writeable. Abort")
		return
	}
	var b = []byte{0x00, 0x00, 0x03, 0x01}

	charToWrite := ch
	for b[3] < 0x8 {
		fmt.Printf("Current value of characteristic : %x\n", b)
		if p.WriteCharacteristic(charToWrite, b, true) != nil {
			fmt.Println("error writing char")
		}

		time.Sleep(5 * time.Second)
		b[3] += 2
	}
	return
}

// onPeriphConnected Callback when a connection to a peripheral is established
func onPeriphConnected(p gatt.Peripheral, err error) {
	fmt.Println("Connected")
	connected <- true
	defer p.Device().CancelConnection(p)
	var numServices int
	var hasErr = false

	// Discover services
	ss, err := p.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return
	}

	svcUUIDNames, _ := csvReadFile("CustomServices.csv")
	charUUIDNames, _ := csvReadFile("CustomCharacteristics.csv")

	xmlDev := &XMLDevice{DeviceName: deviceName}

	for _, s := range ss {
		var svc *XMLService
		var isFoundService bool
		var svcName string
		msg := "Service: " + s.UUID().String()
		if len(s.Name()) > 0 {
			svcName = s.Name()
		} else if len(svcUUIDNames[s.UUID().String()]) > 0 {
			svcName = svcUUIDNames[s.UUID().String()]
		}
		msg += " (" + svcName + ")"
		numServices++
		fmt.Println(msg)
		if isCmpMode == true {
			isFoundService, svc = xmlFindService(device, s.UUID().String())
			if isFoundService == false {
				fmt.Println("Unable to find service ", s.UUID().String(), "in XML Definition")
				hasErr = true
				continue
			}
		}

		// Discover characteristics
		cs, err := p.DiscoverCharacteristics(nil, s)
		if err != nil {
			fmt.Printf("Failed to discover characteristics, err: %s\n", err)
			continue
		}

		var numChars int
		var xmlCharList []XMLCharacteristic
		for _, c := range cs {
			var charName string
			numChars++
			msg := "\tCharacteristic: " + c.UUID().String()
			if len(c.Name()) > 0 {
				charName = c.Name()
			} else if len(charUUIDNames[c.UUID().String()]) > 0 {
				charName = charUUIDNames[c.UUID().String()]
			}
			msg += " (" + charName + ")"
			fmt.Println(msg)
			fmt.Println("\t  ", c.Properties().String())

			ds, err := p.DiscoverDescriptors(nil, c)
			if err != nil {
				fmt.Printf("Failed to discover descriptors, err: %s\n", err)
				continue
			}

			for _, d := range ds {
				msg := "\t\tDescriptor: " + d.UUID().String() + " (" + d.Name() + ") "
				fmt.Println(msg)
			}

			if isCmpMode == true && svc != nil {
				isFoundChar, char := xmlFindChar(svc, c.UUID().String())
				if isFoundChar == false {
					fmt.Println("Unable to find char ", c.UUID().String(), "in XML Definition")
					hasErr = true
					continue
				} else {
					if char.Properties.bitMask != c.Properties() {
						fmt.Println("Char Properties do not match. ")
						fmt.Println("\t Expected '", char.Properties.bitMask,
							"' but found '", c.Properties(), "'")
						hasErr = true
					}
				}
			}
			xmlChar := xmlAppendCharInfo(charName, c.UUID().String(), c.Properties())
			xmlCharList = append(xmlCharList, *xmlChar)
		}
		fmt.Println()
		if (isCmpMode == true) && (numChars != svc.numChars) {
			fmt.Println("Expected", svc.numChars, "characteristics but found", numChars)
			hasErr = true
		}
		xmlSvc := xmlAppendSvcInfo(xmlDev, svcName, s.UUID().String(), xmlCharList)
		xmlDev.ServiceList = append(xmlDev.ServiceList, *xmlSvc)
	}
	if (isCmpMode == true) && (numServices != device.numServices) {
		fmt.Println("Expected", device.numServices, "services but found", numServices)
		hasErr = true
	}

	if isCmpMode == true {
		if hasErr == true {
			fmt.Println("Device did not match specified document")
		} else {
			fmt.Println("Device matches specified document")
		}
	}

	if isXMLMode == true {
		xmlOutDeviceInfo(xmlDev)
	}

	p.Device().CancelConnection(p)
}

// onPeriphDisconnected Callback when a peripheral is disconnected from
func onPeriphDisconnected(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
	close(done)
}

// bleCompareDevice sets up device comparison mode
func bleCompareDevice(macIDArg string, deviceName string, fileName string) {
	isCmpMode = true
	device = xmlGetServices(fileName)

	bleReadDevice(macIDArg, deviceName)
}

// bleReadDeviceXml connects to the specified device and outputs an XML file
func bleReadDeviceXML(macIDArg string, deviceName string) {
	isXMLMode = true
	bleReadDevice(macIDArg, deviceName)
}

// bleReadDevice connects to the specified device
func bleReadDevice(macIDArg string, deviceName string) {
	var err error
	var maxMacLen = 3

	if len(deviceName) == 0 {
		fmt.Println(" Please specify a device to connect to")
		return
	}

	if len(macIDArg) != 0 {
		if len(macIDArg)%2 != 0 {
			fmt.Println("  Invalid len of MAC ID ", macIDArg)
			return
		}
		macID, err = hex.DecodeString(macIDArg)
		if nil != err {
			fmt.Println("  Invalid MAC ID ", macIDArg)
			panic(err)
		}

		if len(macID) != maxMacLen {
			fmt.Println("  Invalid MAC:", macID)
			return
		}
	}

	fmt.Println("\nName: ", deviceName, "\t Identifier: ", macID)

	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("  Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)

	d.Init(onStateChanged)

	bleHandleConnectTimeout()

	<-done
	fmt.Println("Done")
}

// bleScanDevices Scans the radio neighborhood for BLE devices
func bleScanDevices(timeout time.Duration) {
	fmt.Println("Scanning environment for the next", timeout)
	fmt.Println("Please wait ...")

	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	scanMap = make(map[string]ScanMapResult)

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(onScanPeriphDiscovered),
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)

	d.Init(onStateChanged)

	timer1 := time.NewTimer(timeout)

	<-timer1.C

	d.StopScanning()
	displayScanResults()

	bleHandleConnectTimeout()

	<-done
	fmt.Println("Done")
}

func bleHandleConnectTimeout() {
	select {
	case <-connected:

	case <-time.After(maxTimeoutTime):
		fmt.Println("Timed out connecting to device")
		close(done)
	}
}
