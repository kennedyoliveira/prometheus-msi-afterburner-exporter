package monitor

import (
	assertPkg "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_parseXmlResponse(t *testing.T) {
	assert := assertPkg.New(t)

	testFilePath := "test_samples/api_response.xml"

	require.FileExists(t, testFilePath)

	sampleXmlResponse, _ := os.Open(testFilePath)

	parsedResponse, err := parseResponse(sampleXmlResponse)
	if err != nil {
		t.Fatalf("Could not parse xml response: %v", err)
	}

	expectedHeader := HardwareMonitorHeader{
		Signature:     "1296123981",
		Version:       "131072",
		HeaderSize:    32,
		EntryCount:    122,
		EntrySize:     1324,
		Time:          1585285870,
		GpuEntryCount: 1,
		GpuEntrySize:  1304,
	}

	expectedGpuEntry := HardwareMonitorGpuEntry{
		GpuId:     "VEN_10DE&DEV_1B06&SUBSYS_85EA1043&REV_A1&BUS_11&DEV_0&FN_0",
		Family:    "GP102-A",
		Device:    "GeForce GTX 1080 Ti",
		Driver:    "445.75",
		BIOS:      "86.02.39.00.23",
		MemAmount: "0",
	}

	// header
	assert.Equal(expectedHeader, parsedResponse.Header)

	// gpu entries
	assert.Len(parsedResponse.Gpus.Entries, 1)
	assert.Equal(expectedGpuEntry, parsedResponse.Gpus.Entries[0])

	// data metrics

}
