package main

import (
	"fmt"
	"strconv"

	"github.com/mcktr/check_fritz/modules/perfdata"

	"github.com/mcktr/check_fritz/modules/fritz"
	"github.com/mcktr/check_fritz/modules/thresholds"
)

// CheckSmartStatus checks the connection status of a smart device
func CheckSmartStatus(aI ArgumentInformation) {
	soapReq := fritz.NewSoapRequest(*aI.Username, *aI.Password, *aI.Hostname, *aI.Port, "/upnp/control/x_homeauto", "X_AVM-DE_Homeauto", "GetGenericDeviceInfos")
	fritz.AddSoapRequestVariable(&soapReq, fritz.NewSoapRequestVariable("NewIndex", strconv.Itoa(*aI.Index)))

	err := fritz.DoSoapRequest(&soapReq)

	if HandleError(err) {
		return
	}

	var resp = fritz.GetSmartDeviceInfoResponse{}

	err = fritz.HandleSoapRequest(&soapReq, &resp)

	if HandleError(err) {
		return
	}

	output := "- " + resp.NewProductName + " " + resp.NewFirmwareVersion + " - " + resp.NewDeviceName + " " + resp.NewPresent
	GlobalReturnCode = exitOk

	if resp.NewPresent != "CONNECTED" {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK " + output + "\n")
	case exitWarning:
		fmt.Print("WARNING " + output + "\n")
	case exitCritical:
		fmt.Print("CRITICAL " + output + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNWON - Not able to determine smart device status\n")
	}
}

// CheckSmartHeaterTemperatur checks the temperature of a smart home thermometer device
func CheckSmartHeaterTemperatur(aI ArgumentInformation) {
	soapReq := fritz.NewSoapRequest(*aI.Username, *aI.Password, *aI.Hostname, *aI.Port, "/upnp/control/x_homeauto", "X_AVM-DE_Homeauto", "GetGenericDeviceInfos")
	fritz.AddSoapRequestVariable(&soapReq, fritz.NewSoapRequestVariable("NewIndex", strconv.Itoa(*aI.Index)))

	err := fritz.DoSoapRequest(&soapReq)

	if HandleError(err) {
		return
	}

	var resp = fritz.GetSmartDeviceInfoResponse{}

	err = fritz.HandleSoapRequest(&soapReq, &resp)

	if HandleError(err) {
		return
	}

	if resp.NewTemperatureIsEnabled != "ENABLED" {
		fmt.Print("UNKNOWN - Temperature is not enabled on this smart device\n")
		GlobalReturnCode = exitUnknown
		return
	}

	currentTemp, err := strconv.ParseFloat(resp.NewTemperatureCelsius, 64)

	if HandleError(err) {
		return
	}

	currentTemp = currentTemp / 10.0
	perfData := perfdata.CreatePerformanceData("temperature", currentTemp, "")

	GlobalReturnCode = exitOk

	if thresholds.IsSet(aI.Warning) {
		perfData.SetWarning(*aI.Warning)

		if thresholds.CheckLower(*aI.Warning, currentTemp) {
			GlobalReturnCode = exitWarning
		}
	}

	if thresholds.IsSet(aI.Critical) {
		perfData.SetCritical(*aI.Critical)

		if thresholds.CheckLower(*aI.Critical, currentTemp) {
			GlobalReturnCode = exitCritical
		}
	}

	output := "- " + resp.NewProductName + " " + resp.NewFirmwareVersion + " - " + resp.NewDeviceName + " " + fmt.Sprintf("%.2f", currentTemp) + " °C " + perfData.GetPerformanceDataAsString()

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK " + output + "\n")
	case exitWarning:
		fmt.Print("WARNING " + output + "\n")
	case exitCritical:
		fmt.Print("CRITICAL " + output + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNWON - Not able to calculate heater temperature\n")
	}
}

// CheckSmartSocketPower checks the current watt usage on the smart socket
func CheckSmartSocketPower(aI ArgumentInformation) {
	soapReq := fritz.NewSoapRequest(*aI.Username, *aI.Password, *aI.Hostname, *aI.Port, "/upnp/control/x_homeauto", "X_AVM-DE_Homeauto", "GetGenericDeviceInfos")
	fritz.AddSoapRequestVariable(&soapReq, fritz.NewSoapRequestVariable("NewIndex", strconv.Itoa(*aI.Index)))

	err := fritz.DoSoapRequest(&soapReq)

	if HandleError(err) {
		return
	}

	var resp = fritz.GetSmartDeviceInfoResponse{}

	err = fritz.HandleSoapRequest(&soapReq, &resp)

	if HandleError(err) {
		return
	}

	currentPower, err := strconv.ParseFloat(resp.NewMultimeterPower, 64)

	if HandleError(err) {
		return
	}

	currentPower = currentPower / 100.0
	perfData := perfdata.CreatePerformanceData("power", currentPower, "")

	GlobalReturnCode = exitOk

	if thresholds.IsSet(aI.Warning) {
		perfData.SetWarning(*aI.Warning)

		if thresholds.CheckUpper(*aI.Warning, currentPower) {
			GlobalReturnCode = exitWarning
		}
	}

	if thresholds.IsSet(aI.Critical) {
		perfData.SetCritical(*aI.Critical)

		if thresholds.CheckUpper(*aI.Critical, currentPower) {
			GlobalReturnCode = exitCritical
		}
	}

	output := "- " + resp.NewProductName + " " + resp.NewFirmwareVersion + " - " + resp.NewDeviceName + " " + fmt.Sprintf("%.2f", currentPower) + " W " + perfData.GetPerformanceDataAsString()

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK " + output + "\n")
	case exitWarning:
		fmt.Print("WARNING " + output + "\n")
	case exitCritical:
		fmt.Print("CRITICAL " + output + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNWON - Not able to fetch socket power\n")
	}
}

// CheckSmartSocketEnergy checks total power consumption of the last year on the smart socket
func CheckSmartSocketEnergy(aI ArgumentInformation) {
	soapReq := fritz.NewSoapRequest(*aI.Username, *aI.Password, *aI.Hostname, *aI.Port, "/upnp/control/x_homeauto", "X_AVM-DE_Homeauto", "GetGenericDeviceInfos")
	fritz.AddSoapRequestVariable(&soapReq, fritz.NewSoapRequestVariable("NewIndex", strconv.Itoa(*aI.Index)))

	err := fritz.DoSoapRequest(&soapReq)

	if HandleError(err) {
		return
	}

	var resp = fritz.GetSmartDeviceInfoResponse{}

	err = fritz.HandleSoapRequest(&soapReq, &resp)

	if HandleError(err) {
		return
	}

	currentEnergy, err := strconv.ParseFloat(resp.NewMultimeterEnergy, 64)

	if HandleError(err) {
		return
	}

	currentEnergy = currentEnergy / 1000.0
	perfData := perfdata.CreatePerformanceData("energy", currentEnergy, "")

	GlobalReturnCode = exitOk

	if thresholds.IsSet(aI.Warning) {
		perfData.SetWarning(*aI.Warning)

		if thresholds.CheckUpper(*aI.Warning, currentEnergy) {
			GlobalReturnCode = exitWarning
		}
	}

	if thresholds.IsSet(aI.Critical) {
		perfData.SetCritical(*aI.Critical)

		if thresholds.CheckUpper(*aI.Critical, currentEnergy) {
			GlobalReturnCode = exitCritical
		}
	}

	output := "- " + resp.NewProductName + " " + resp.NewFirmwareVersion + " - " + resp.NewDeviceName + " " + fmt.Sprintf("%.2f", currentEnergy) + " kWh " + perfData.GetPerformanceDataAsString()

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK " + output + "\n")
	case exitWarning:
		fmt.Print("WARNING " + output + "\n")
	case exitCritical:
		fmt.Print("CRITICAL " + output + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNWON - Not able to fetch socket energy\n")
	}
}
