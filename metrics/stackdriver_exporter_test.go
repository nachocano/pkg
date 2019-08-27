/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"knative.dev/pkg/metrics/monitoredresources"
	"path"
	"testing"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	. "knative.dev/pkg/logging/testing"
	"knative.dev/pkg/metrics/metricskey"
	metricskeyeventing "knative.dev/pkg/metrics/metricskey/eventing"
	metricskeyserving "knative.dev/pkg/metrics/metricskey/serving"
)

// TODO should be properly refactored and pieces should move to eventing and serving, as appropriate.
// 	See https://github.com/knative/pkg/issues/608

var (
	testGcpMetadata = monitoredresources.GcpMetadata{
		Project:  "test-project",
		Location: "test-location",
		Cluster:  "test-cluster",
	}

	supportedServingMetricsTestCases = []struct {
		name       string
		domain     string
		component  string
		metricName string
	}{{
		name:       "activator metric",
		domain:     servingDomain,
		component:  "activator",
		metricName: "request_count",
	}, {
		name:       "autoscaler metric",
		domain:     servingDomain,
		component:  "autoscaler",
		metricName: "desired_pods",
	}}

	supportedEventingBrokerMetricsTestCases = []struct {
		name       string
		domain     string
		component  string
		metricName string
	}{{
		name:       "broker metric",
		domain:     eventingDomain,
		component:  "broker",
		metricName: "event_count",
	}}

	supportedEventingTriggerMetricsTestCases = []struct {
		name       string
		domain     string
		component  string
		metricName string
	}{{
		name:       "trigger metric",
		domain:     eventingDomain,
		component:  "trigger",
		metricName: "event_count",
	}, {
		name:       "trigger metric",
		domain:     eventingDomain,
		component:  "trigger",
		metricName: "event_process_latencies",
	}}

	unsupportedMetricsTestCases = []struct {
		name       string
		domain     string
		component  string
		metricName string
	}{{
		name:       "unsupported domain",
		domain:     "unsupported",
		component:  "activator",
		metricName: "request_count",
	}, {
		name:       "unsupported component",
		domain:     servingDomain,
		component:  "unsupported",
		metricName: "request_count",
	}, {
		name:       "unsupported metric",
		domain:     servingDomain,
		component:  "activator",
		metricName: "unsupported",
	}, {
		name:       "unsupported component",
		domain:     eventingDomain,
		component:  "unsupported",
		metricName: "event_count",
	}, {
		name:       "unsupported metric",
		domain:     eventingDomain,
		component:  "broker",
		metricName: "unsupported",
	}}
)

func fakeGcpMetadataFun() *monitoredresources.GcpMetadata {
	return &testGcpMetadata
}

type fakeExporter struct{}

func (fe *fakeExporter) ExportView(vd *view.Data) {}
func (fe *fakeExporter) Flush()                   {}

func newFakeExporter(o stackdriver.Options) (view.Exporter, error) {
	return &fakeExporter{}, nil
}

func TestGetMonitoredResourceFunc_UseKnativeRevision(t *testing.T) {
	for _, testCase := range supportedServingMetricsTestCases {
		testView = &view.View{
			Description: "Test View",
			Measure:     stats.Int64(testCase.metricName, "Test Measure", stats.UnitNone),
			Aggregation: view.LastValue(),
			TagKeys:     []tag.Key{},
		}
		mrf := getMonitoredResourceFunc(path.Join(testCase.domain, testCase.component), &testGcpMetadata)

		newTags, monitoredResource := mrf(testView, testTags)
		gotResType, labels := monitoredResource.MonitoredResource()
		wantedResType := "knative_revision"
		if gotResType != wantedResType {
			t.Fatalf("MonitoredResource=%v, want %v", gotResType, wantedResType)
		}
		got := getResourceLabelValue(metricskeyserving.LabelRouteName, newTags)
		if got != testRoute {
			t.Errorf("expected new tag: %v, got: %v", routeKey, newTags)
		}
		got, ok := labels[metricskey.LabelNamespaceName]
		if !ok || got != testNS {
			t.Errorf("expected label %v with value %v, got: %v", metricskey.LabelNamespaceName, testNS, got)
		}
		got, ok = labels[metricskeyserving.LabelConfigurationName]
		if !ok || got != metricskey.ValueUnknown {
			t.Errorf("expected label %v with value %v, got: %v", metricskeyserving.LabelConfigurationName, metricskey.ValueUnknown, got)
		}
	}
}

func TestGetMonitoredResourceFunc_UseKnativeBroker(t *testing.T) {
	for _, testCase := range supportedEventingBrokerMetricsTestCases {
		testView = &view.View{
			Description: "Test View",
			Measure:     stats.Int64(testCase.metricName, "Test Measure", stats.UnitNone),
			Aggregation: view.LastValue(),
			TagKeys:     []tag.Key{},
		}
		mrf := getMonitoredResourceFunc(path.Join(testCase.domain, testCase.component), &testGcpMetadata)

		newTags, monitoredResource := mrf(testView, testTags)
		gotResType, labels := monitoredResource.MonitoredResource()
		wantedResType := "knative_broker"
		if gotResType != wantedResType {
			t.Fatalf("MonitoredResource=%v, want %v", gotResType, wantedResType)
		}
		got := getResourceLabelValue(metricskeyeventing.LabelBrokerName, newTags)
		if got != testBroker {
			t.Errorf("expected new tag: %v, got: %v", brokerKey, newTags)
		}
		got, ok := labels[metricskey.LabelNamespaceName]
		if !ok || got != testNS {
			t.Errorf("expected label %v with value %v, got: %v", metricskey.LabelNamespaceName, testNS, got)
		}
	}
}

func TestGetMonitoredResourceFunc_UseGlobal(t *testing.T) {
	for _, testCase := range unsupportedMetricsTestCases {
		testView = &view.View{
			Description: "Test View",
			Measure:     stats.Int64(testCase.metricName, "Test Measure", stats.UnitNone),
			Aggregation: view.LastValue(),
			TagKeys:     []tag.Key{},
		}
		mrf := getMonitoredResourceFunc(path.Join(testCase.domain, testCase.component), &testGcpMetadata)

		newTags, monitoredResource := mrf(testView, testTags)
		gotResType, labels := monitoredResource.MonitoredResource()
		wantedResType := "global"
		if gotResType != wantedResType {
			t.Fatalf("MonitoredResource=%v, want: %v", gotResType, wantedResType)
		}
		got := getResourceLabelValue(metricskey.LabelNamespaceName, newTags)
		if got != testNS {
			t.Errorf("expected new tag %v with value %v, got: %v", routeKey, testNS, newTags)
		}
		if len(labels) != 0 {
			t.Errorf("expected no label, got: %v", labels)
		}
	}
}

func TestGetgetMetricTypeFunc_UseKnativeDomain(t *testing.T) {
	for _, testCase := range supportedServingMetricsTestCases {
		testView = &view.View{
			Description: "Test View",
			Measure:     stats.Int64(testCase.metricName, "Test Measure", stats.UnitNone),
			Aggregation: view.LastValue(),
			TagKeys:     []tag.Key{},
		}
		mtf := getMetricTypeFunc(
			path.Join(testCase.domain, testCase.component),
			path.Join(customMetricTypePrefix, testCase.component))

		gotMetricType := mtf(testView)
		wantedMetricType := path.Join(testCase.domain, testCase.component, testView.Measure.Name())
		if gotMetricType != wantedMetricType {
			t.Fatalf("getMetricType=%v, want %v", gotMetricType, wantedMetricType)
		}
	}
}

func TestGetgetMetricTypeFunc_UseCustomDomain(t *testing.T) {
	for _, testCase := range unsupportedMetricsTestCases {
		testView = &view.View{
			Description: "Test View",
			Measure:     stats.Int64(testCase.metricName, "Test Measure", stats.UnitNone),
			Aggregation: view.LastValue(),
			TagKeys:     []tag.Key{},
		}
		mtf := getMetricTypeFunc(
			path.Join(testCase.domain, testCase.component),
			path.Join(customMetricTypePrefix, testCase.component))

		gotMetricType := mtf(testView)
		wantedMetricType := path.Join(customMetricTypePrefix, testCase.component, testView.Measure.Name())
		if gotMetricType != wantedMetricType {
			t.Fatalf("getMetricType=%v, want %v", gotMetricType, wantedMetricType)
		}
	}
}

func TestNewStackdriverExporterWithMetadata(t *testing.T) {
	e, err := newStackdriverExporter(&metricsConfig{
		domain:               servingDomain,
		component:            "autoscaler",
		backendDestination:   Stackdriver,
		stackdriverProjectID: testProj}, TestLogger(t))
	if err != nil {
		t.Error(err)
	}
	if e == nil {
		t.Error("expected a non-nil metrics exporter")
	}
}
