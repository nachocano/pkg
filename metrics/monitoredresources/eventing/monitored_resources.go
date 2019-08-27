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

package eventing

// TODO should be moved to eventing. See https://github.com/knative/pkg/issues/608

import (
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"knative.dev/pkg/metrics/metricskey"
	metricskeyeventing "knative.dev/pkg/metrics/metricskey/eventing"
	"knative.dev/pkg/metrics/monitoredresources"
)

type KnativeTrigger struct {
	Project               string
	Location              string
	ClusterName           string
	NamespaceName         string
	TriggerName           string
	BrokerName            string
	TypeFilterAttribute   string
	SourceFilterAttribute string
}

type KnativeBroker struct {
	Project       string
	Location      string
	ClusterName   string
	NamespaceName string
	BrokerName    string
}

type KnativeImporter struct {
	Project       string
	Location      string
	ClusterName   string
	NamespaceName string
	ImporterName  string
	ImporterKind  string
}

func (kt *KnativeTrigger) MonitoredResource() (resType string, labels map[string]string) {
	labels = map[string]string{
		metricskey.LabelProject:                              kt.Project,
		metricskey.LabelLocation:                             kt.Location,
		metricskey.LabelClusterName:                          kt.ClusterName,
		metricskey.LabelNamespaceName:                        kt.NamespaceName,
		metricskeyeventing.LabelTriggerName:                  kt.TriggerName,
		metricskeyeventing.LabelBrokerName:                   kt.BrokerName,
		metricskeyeventing.LabelTriggerTypeFilterAttribute:   kt.TypeFilterAttribute,
		metricskeyeventing.LabelTriggerSourceFilterAttribute: kt.SourceFilterAttribute,
	}
	return "knative_trigger", labels
}

func (kb *KnativeBroker) MonitoredResource() (resType string, labels map[string]string) {
	labels = map[string]string{
		metricskey.LabelProject:            kb.Project,
		metricskey.LabelLocation:           kb.Location,
		metricskey.LabelClusterName:        kb.ClusterName,
		metricskey.LabelNamespaceName:      kb.NamespaceName,
		metricskeyeventing.LabelBrokerName: kb.BrokerName,
	}
	return "knative_broker", labels
}

func (ki *KnativeImporter) MonitoredResource() (resType string, labels map[string]string) {
	labels = map[string]string{
		metricskey.LabelProject:              ki.Project,
		metricskey.LabelLocation:             ki.Location,
		metricskey.LabelClusterName:          ki.ClusterName,
		metricskey.LabelNamespaceName:        ki.NamespaceName,
		metricskeyeventing.LabelImporterName: ki.ImporterName,
		metricskeyeventing.LabelImporterKind: ki.ImporterKind,
	}
	return "knative_importer", labels
}

func GetKnativeBrokerMonitoredResource(
	v *view.View, tags []tag.Tag, gm *monitoredresources.GcpMetadata) ([]tag.Tag, monitoredresource.Interface) {
	tagsMap := monitoredresources.GetTagsMap(tags)
	kb := &KnativeBroker{
		// The first three resource labels are from metadata.
		Project:     gm.Project,
		Location:    gm.Location,
		ClusterName: gm.Cluster,
		// The rest resource labels are from metrics labels.
		NamespaceName: monitoredresources.ValueOrUnknown(metricskey.LabelNamespaceName, tagsMap),
		BrokerName:    monitoredresources.ValueOrUnknown(metricskeyeventing.LabelBrokerName, tagsMap),
	}

	var newTags []tag.Tag
	for _, t := range tags {
		// Keep the metrics labels that are not resource labels
		if !metricskeyeventing.KnativeBrokerLabels.Has(t.Key.Name()) {
			newTags = append(newTags, t)
		}
	}

	return newTags, kb
}

func GetKnativeTriggerMonitoredResource(
	v *view.View, tags []tag.Tag, gm *monitoredresources.GcpMetadata) ([]tag.Tag, monitoredresource.Interface) {
	tagsMap := monitoredresources.GetTagsMap(tags)
	kt := &KnativeTrigger{
		// The first three resource labels are from metadata.
		Project:     gm.Project,
		Location:    gm.Location,
		ClusterName: gm.Cluster,
		// The rest resource labels are from metrics labels.
		NamespaceName:         monitoredresources.ValueOrUnknown(metricskey.LabelNamespaceName, tagsMap),
		TriggerName:           monitoredresources.ValueOrUnknown(metricskeyeventing.LabelTriggerName, tagsMap),
		BrokerName:            monitoredresources.ValueOrUnknown(metricskeyeventing.LabelBrokerName, tagsMap),
		TypeFilterAttribute:   monitoredresources.ValueOrUnknown(metricskeyeventing.LabelTriggerTypeFilterAttribute, tagsMap),
		SourceFilterAttribute: monitoredresources.ValueOrUnknown(metricskeyeventing.LabelTriggerSourceFilterAttribute, tagsMap),
	}

	var newTags []tag.Tag
	for _, t := range tags {
		// Keep the metrics labels that are not resource labels
		if !metricskeyeventing.KnativeTriggerLabels.Has(t.Key.Name()) {
			newTags = append(newTags, t)
		}
	}

	return newTags, kt
}

func GetKnativeImporterMonitoredResource(
	v *view.View, tags []tag.Tag, gm *monitoredresources.GcpMetadata) ([]tag.Tag, monitoredresource.Interface) {
	tagsMap := monitoredresources.GetTagsMap(tags)
	ki := &KnativeImporter{
		// The first three resource labels are from metadata.
		Project:     gm.Project,
		Location:    gm.Location,
		ClusterName: gm.Cluster,
		// The rest resource labels are from metrics labels.
		NamespaceName: monitoredresources.ValueOrUnknown(metricskey.LabelNamespaceName, tagsMap),
		ImporterName:  monitoredresources.ValueOrUnknown(metricskeyeventing.LabelImporterName, tagsMap),
		ImporterKind:  monitoredresources.ValueOrUnknown(metricskeyeventing.LabelImporterKind, tagsMap),
	}

	var newTags []tag.Tag
	for _, t := range tags {
		// Keep the metrics labels that are not resource labels
		if !metricskeyeventing.KnativeImporterLabels.Has(t.Key.Name()) {
			newTags = append(newTags, t)
		}
	}

	return newTags, ki
}
