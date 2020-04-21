package monitor

import (
	"testing"
)

// Should transform metric names from MSIAfterBurner
// in valid metrics for prometheus
func Test_normalizeMetricName(t *testing.T) {

	onlyName := func(name string) *HardwareMonitorEntry {
		return &HardwareMonitorEntry{LocalizedSourceName: name}
	}

	tests := []struct {
		name   string
		metric *HardwareMonitorEntry
		want   string
	}{
		{name: "Valid String", metric: onlyName("gpu_temperature_celsius"), want: "gpu_temperature_celsius"},
		{name: "Upper case should be lowered", metric: onlyName("GPU_TEMPERATURE_CELSIUS"), want: "gpu_temperature_celsius"},
		{name: "Convert space to _", metric: onlyName("Gpu Temperature"), want: "gpu_temperature"},
		{name: "Multiple space to single _", metric: onlyName("Gpu    Temperature"), want: "gpu_temperature"},
		{name: "Leading spaces must be removed", metric: onlyName("  Gpu  Temperature"), want: "gpu_temperature"},
		{name: "Non valid characters should be removed", metric: onlyName(" Gpu 123 temperature "), want: "gpu_temperature"},
		{name: "Remove invalid characters from the beginning of the string", metric: onlyName("123gpu_temperature"), want: "gpu_temperature"},
		{name: "CPU Metrics should extract cpu core as label", metric: onlyName("CPU18 Temperature"), want: "cpu_temperature"},
		{name: "Metrics of unknown type should have no suffix", metric: onlyName("Voltage limit"), want: "voltage_limit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := normalizeMetricName(tt.metric, nil); got != tt.want {
				t.Errorf("normalizeMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeMetric(t *testing.T) {
	type args struct {
		metricValue float64
		metricUnit  string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "Percentage metric should be converted to 0 ~ 1", args: args{metricValue: 53.5, metricUnit: "%"}, want: 0.535},
		{name: "Convert MB to bytes", args: args{metricValue: 35.3, metricUnit: "MB"}, want: 3.53e+7},
		{name: "Convert mega hertz to hertz", args: args{metricValue: 17.3, metricUnit: "MHz"}, want: 1.73e+7},
		{name: "Other values should be just returned", args: args{metricValue: 25.32, metricUnit: "Some other stuff"}, want: 25.32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeMetricValue(tt.args.metricValue, tt.args.metricUnit); got != tt.want {
				t.Errorf("normalizeMetricValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
