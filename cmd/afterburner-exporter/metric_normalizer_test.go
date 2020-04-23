package main

import (
	"github.com/kennedyoliveira/prometheus-msi-afterburner-exporter/afterburner"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
	"testing"
)

func Test_normalizeMetricName(t *testing.T) {
	type args struct {
		metric *afterburner.HardwareMonitorEntry
		gpus   *[]afterburner.HardwareMonitorGpuEntry
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   *prometheus.Labels
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := normalizeMetricName(tt.args.metric, tt.args.gpus)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeMetricName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeMetricName() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("normalizeMetricName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
