---
layout: post
title: Introducing BLE Tools
---

# BLE Tools

## Bluetooth Low Energy (BLE)
Bluetooth is a wireless personal area networking standard for exchanging data
over short distances.  Bluetooth low energy (BLE) (also known as Version 4.0+
of the Bluetooth specification, or Bluetooth Smart) is the power- and
application-friendly version of Bluetooth that was built for the Internet of
Things (IoT). The power efficiency and low energy functionality make this protocol
perfect for battery operated devices. Since BLE now comes native on every modern
phone, tablet and computer, it also makes for a perfect starting point to connect
with the vast multitude of devices that IoT promises to bring to the world.

At [Blue Clover Devices](http://www.bcdevices.com) (BCD),  we work on a wide
variety of IoT products, with BLE providing connectivity for a significant portion of them.
One such example is a BLE controlled desk lamp named
[Sensalite](https://itunes.apple.com/us/app/sensalite/id1053228450?mt=8).
Using an iOS app, one can toggle the light switch, as well as
adjust the brightness of the lamp. The app also provides additional information
such as the battery level, and if the device is powered by USB or running on battery.
The firmware and the app communicate with each other via BLE.

## Developing with BLE
An IoT project that use BLE for connectivity consists of several components, namely

* The device: which includes design and implementation which includes the hardware and firmware
* The apps: which includes the overall user experience and app development for different platforms
* The cloud backend
* The BLE interface between a mobile app and the device

The focus of this discussion is this last element.
Design of the BLE interface includes two aspects:

1. Defining the BLE advertisement and scan response content
1. Designing the BLE profile


### Advertisements
BLE advertisements are a periodic unidirectional broadcast from the peripheral to all
devices around it. A central can then use the information in these packets to
connect to the peripheral. The advertising interval can be fixed to be between
20ms and 10.24s, in steps of 0.625ms. A random delay of up to 10ms is also added
to the fixed interval to reduce the possibility of collisions between advertisements
of different peripherals.

The most important aspect of building the advertisement is providing the relevant
information to a device that wants to connect. Optionally, information that uniquely
identifies the peripheral may also be included. If privacy is a concern, this
information may be skipped. The advertisement has 31 bytes that can be used to
advertise different things. The most common advertisement payloads are:
*  Incomplete list of 128 bit services
*  Incomplete list of 16-bit service class UUID
*  Local Name
*  Manufacturer Specific Data
*  Power Level

### GATT (Generic ATTribute Profile)
The GATT is a set of rules describing how to bundle and present data
using BLE. A GATT consists of Services and Characteristics.

The Bluetooth Core Specification defines the GATT as:

> The GATT Profile specifies
> the structure in which profile data is exchanged. This structure defines basic
> elements such as services and characteristics, used in a profile.

A good visualization of this is to imagine a bookshelf with each shelf containing
books of different topics. The bookshelf would be your GATT profile; every shelf
would be a service and each book would be a characteristic.

These services and characteristics are identified by Universally Unique IDs (UUID).
The Bluetooth SIG predefines several UUIDs for the standard services, and these
are 16-bit UUIDs. An example is the Battery Service, which has a UUID of `0x180F`.
This service will contain the Battery Level Characteristic, which is defined to
have a UUID of `0x2A19`. When defining custom services and characteristics, one
would use a 128-bit UUID for which an example would look something like
`1C68B3FA-D443-4365-9E1C-B22F44EB0816`.

#### Services
A Service is defined as: "_[...] a collection of data and associated behaviors
to accomplish a particular function or feature._"
Specifically, a service is a logical collection of characteristics. A service can
have one or more characteristics.

#### Characteristics
A Characteristic is defined as, "_[...] a value used in a service along with
properties and configuration information about how the value is accessed and
information about how the value is displayed or represented._"
As such, a characteristic is the true presentation of the information being
represented. Characteristics can have a combination of several properties.
Some of the properties that a characteristic can possess are:

* Broadcast - allows the characteristic to be placed in advertisements
* Read - allows clients to read the value of this characteristic using any Read ATT operation
* Write without Response - allows clients to use the Write Command ATT operation
* Write - allows clients to use the Write Request/Response ATT operation
* Notify - allows the server to send unacked data to the client
* Indicate - allows the server to send acknowledged data to the client

## ble-tools
We experienced a lack of command line tools for testing and analysis as we developed
different BLE products. This led us to build our own. The benefit of building a command line
tool is its expandability for scripting, and building automation to help test any
product being built.  For a GUI based tool, take a look at [Light Blue Explorer](https://itunes.apple.com/us/app/lightblue-explorer-bluetooth/id557428110?mt=8).

With the basic understanding of BLE above, we can now look at the tool we have
developed in some detail. Some of the things we aimed to cover in BLE development
were checking the advertisements packets, being able to connect to a peripheral,
and also extract the GATT definitions to compare between the implementation and the
design.

### Usage
The tool supports 4 basic modes:

1. Scan for devices
1. Connect to specific device and optionally generate an XML-based GATT Definition
1. Parse XML input file and generate human readable output of a Device's GATT Definition
1. Compare Device Profile with the input XML definition

### Github
The project can be found on [Github](https://github.com/bcdevices/ble-tools).

Ensure that it is cloned into `$GOPATH/src/github.com/bcdevices/ble-tools`

Run `make clean && make` to end up with the relevant artifacts in the `$GOPATH/bin` folder.

## Future Work
We consider this to only be the first stage of development of this tool. Some things
we are still in the process of building are:

1. Testing for poorly handled write handlers, including
  * Sending bytes larger than the specified size
  * Sending data outside of the valid range specified in the interface
1. Testing the connectivity of the device
  * Running a loop to rapidly connect and disconnect from the peripheral

## References

1. https://www.bluetooth.org/
1. https://devzone.nordicsemi.com/tutorials/8/ble-services-a-beginners-tutorial/
1. http://www.argenox.com/a-ble-advertising-primer/
