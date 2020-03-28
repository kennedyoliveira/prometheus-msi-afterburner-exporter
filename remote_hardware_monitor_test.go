package main

import "testing"

// Should transform metric names from MSIAfterBurner
// in valid metrics for prometheus
func Test_normalizeMetricName(t *testing.T) {
	tests := []struct {
		name       string
		metricName string
		want       string
	}{
		{name: "Valid String", metricName: "gpu_temperature_celsius", want: "gpu_temperature_celsius"},
		{name: "Upper case should be lowered", metricName: "GPU_TEMPERATURE_CELSIUS", want: "gpu_temperature_celsius"},
		{name: "Convert space to _", metricName: "Gpu Temperature", want: "gpu_temperature"},
		{name: "Multiple space to single _", metricName: "Gpu    Temperature", want: "gpu_temperature"},
		{name: "Leading spaces must be removed", metricName: "  Gpu  Temperature", want: "gpu_temperature"},
		{name: "Non valid characters should be removed", metricName: " Gpu 123 temperature ", want: "gpu_temperature"},
		{name: "Remove invalid characters from the begnning of the string", metricName: "123gpu_temperature", want: "gpu_temperature"},
		{name: "CPU Metrics should extract cpu core as label", metricName: "CPU18 Temperature", want: "cpu_temperature"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := normalizeMetricName(tt.metricName, ""); got != tt.want {
				t.Errorf("normalizeMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}
