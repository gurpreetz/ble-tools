# ble-tools

[![Build Status](https://drone.io/github.com/bcdevices/ble-tools/status.png)](https://drone.io/github.com/bcdevices/ble-tools/latest)
[![GoDoc](https://godoc.org/github.com/bcdevices/ble-tools?status.svg)](https://godoc.org/github.com/bcdevices/ble-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/bcdevices/ble-tools)](https://goreportcard.com/report/github.com/bcdevices/ble-tools)


Software tool to test BLE Devices

- Build platform: OS X
- Host platform: OS X


## Usage
The device supports 4 basic modes:

1. Scan for devices
1. Connect to specific device
1. Read XML input file that defines a device
1. Compare Physical Device with XML definitions

The basic modes of usage for ble-tools can be seen below:

    Usage: ./ble-tools [COMMAND] [<options>]
    scan
      -timeout timeout
        	scan timeout duration in seconds (default 12s)
    connect
      -device Device Name
        	BLE Device Name
      -id mfg data
        	Last 3 hex bytes of mfg data to uniquely identify device
      -xmlOut
        	generate an xml output
    read
      -file xml file
        	xml file to be parsed
    compare
      -device Device Name
        	BLE Device Name
      -file XML file
        	XML file to compare against
      -id mfg data
        	Last 3 hex bytes of mfg data to uniquely identify device

### Scan
This runs a passive scan of the neighboring environment for the duration of time specified
with the `timeout` flag. By default, this timeout is 12s. This is the maximum advertising 
timeout specified in the BLE Specification. The minimum allowed timeout value is 1s. 
After a scan is completed, a numbered list of discovered devices is displayed. The user can select
the device to connect to by entering the index of the device. 
An example of a scan list is shown below. 

    Following Devices discovered:
        Index    Device Name
         0   estimote-2892c4
         1   Apple TV-060090
         2   T2-000017
         3   Unknown-d1c887
         4   Dropcam-0ff4c7
         5   Unknown-3d05a1
         6   estimote-91c4c4
         7   Aug-d10100
         8   BCD Sensalite-000094
    Enter Device to connect to: 2
    Generate XML after discovery? (y/n) [n]:

Upon making a selection, the user is prompted whether they would like an XML output of the
device's services and characteristics to be saved. If desired, this is generated after connecting
to the device, and saved in the `XmlOutputs` folder. 

### Connect
Many devices announce their names in the LocalName field of the BLE advertisement. If one already
knows this name, and would like to connect to the device without having to explicitly scan the 
environment, this would be the preferred mode to use. The name of the device is specified in the 
mandatory field of `device`. For further fine-grained connection options, one can also specify the
last three bytes of the manufacturing data from the advertisement in the optional `id` field. 
Finally, if one would like a record of the device's services and characteristics
a true/false field of `xmlOut` can be used as well. By default this option is set to false. 

#### Custom Services/Characteristics
Not all devices use the standard Bluetooth specified service/characteristic UUIDs. To help in making
this information readable two user generated files are included, viz `CustomServices.csv`
and `CustomCharacteristics.csv`. 
These are CSV files that have a mapping of UUID to a Human Readable Service or Characteristic name,
respectively. Currently this only contains mappings for devices built here at BCD, but it can grow
with time as this tool hopefully gets used. Another use would be populating it with the Apple UUIDs 
defined in the HomeKit Specification. 

### Read
Once an xml file of the device's services and characteristics  has already been generated, either by 
this tool, or by other means, this tool can parse the information in the file and display it in a human 
readable format. The file to be read is passed in via the mandatory option of `file`. 

An XML file representing the device consists of one or many services, each of which can contain several
characteristics. The properties of each characteristic can be obtained from the Bluetooth website. 
A snippet of one device would look like this. 

    <?xml  version="1.0" encoding="UTF-8"?>
    <device name="Ly01">
    <service name="LY-01 Service" uuid="1c68b3fad44343659e1cb22f44eb0816">
        <characteristic name="Light Control" uuid="447c291d5318420b980a8f33e22c3744">
            <requirement>Mandatory</requirement>
            <Properties>
                <Read>Mandatory</Read>
                <Write>Mandatory</Write>
                <WriteWithoutResponse>Excluded</WriteWithoutResponse>
                <SignedWrite>Excluded</SignedWrite>
                <ReliableWrite>Excluded</ReliableWrite>
                <Notify>Mandatory</Notify>
                <Indicate>Mandatory</Indicate>
                <WritableAuxilliaries>Excluded</WritableAuxilliaries>
                <Broadcast>Excluded</Broadcast>
            </Properties>
        </characteristic>
    </service> 
    </device>

### Compare
While building a device, its always useful to ensure that the BLE interface on the device
matches what was specified in the interface design document. The `compare` mode helps achieve this goal. 
This mode is a combination of the `read` and `connect` modes described above. The tool will read the file
specified, and connect to the device in question. It will then scan the device, and perform  a property by
property comparison across all characteristics and services. All inconsistencies detected will be reported 
after disconnection from the device. 
An example of some errors that can be reported are:
   
    Char Properties do not match. 
         Expected ' read  ' but found ' read notify  '
    
    Expected 9 characteristics but found 8
    
    Expected 4 services but found 3
    
    Device did not match specified document

## Local build

- Ensure the repository is checked out in `$GOPATH/src/github.com/bcdevices/ble-tools`
- In toplevel directory, execute:

        make clean
        make

to end up with the build artifacts in the `$GOPATH/bin` folder.

### Prerequisites

- Apple Mac computer running OS X
- OS X: Xcode Command Line Tools
- Go 1.6
