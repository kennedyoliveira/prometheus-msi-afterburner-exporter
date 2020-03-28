package main

import (
	"encoding/xml"
	"io"
)

type HardwareMonitor struct {
	Header HardwareMonitorHeader     `xml:"HardwareMonitorHeader"`
	Data   HardwareMonitorEntries    `xml:"HardwareMonitorEntries"`
	Gpus   HardwareMonitorGpuEntries `xml:"HardwareMonitorGpuEntries"`
}

type HardwareMonitorHeader struct {
	Signature     string `xml:"signature"`
	Version       string `xml:"version"`
	HeaderSize    int    `xml:"headerSize"`
	EntryCount    int    `xml:"entryCount"`
	EntrySize     int    `xml:"entrySize"`
	Time          int64  `xml:"time"`
	GpuEntryCount int    `xml:"gpuEntryCount"`
	GpuEntrySize  int    `xml:"gpuEntrySize"`
}

type HardwareMonitorGpuEntries struct {
	Entries []HardwareMonitorGpuEntry `xml:"HardwareMonitorGpuEntry"`
}

type HardwareMonitorGpuEntry struct {
	GpuId     string `xml:"gpuId"`
	Family    string `xml:"family"`
	Device    string `xml:"device"`
	Driver    string `xml:"driver"`
	BIOS      string `xml:"BIOS"`
	MemAmount string `xml:"memAmount"`
}

type HardwareMonitorEntries struct {
	Entries []HardwareMonitorEntry `xml:"HardwareMonitorEntry"`
}

type HardwareMonitorEntry struct {
	SourceName          string  `xml:"srcName"`
	SourceUnits         string  `xml:"srcUnit"`
	LocalizedSourceName string  `xml:"localizedSrcName"`
	LocalizedSourceUnit string  `xml:"localizedSrcUnits"`
	RecommendedFormat   string  `xml:"recommendedFormat"`
	Data                float64 `xml:"data"`
	MinLimit            float64 `xml:"minLimit"`
	MaxLimit            float64 `xml:"maxLimit"`
	Flags               string  `xml:"flags"`
	GpuIndex            int     `xml:"gpu"`
	SourceId            string  `xml:"srcId"`
}

func parseResponse(reader io.Reader) (*HardwareMonitor, error) {
	var hardwareMonitor HardwareMonitor

	err := xml.NewDecoder(reader).Decode(&hardwareMonitor)
	if err != nil {
		return nil, err
	}

	return &hardwareMonitor, nil
}
