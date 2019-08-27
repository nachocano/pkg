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

package serving

import (
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"knative.dev/pkg/metrics/metricskey"
	metricskeyserving "knative.dev/pkg/metrics/metricskey/serving"
	"knative.dev/pkg/metrics/monitoredresources"
)

// TODO should be moved to serving. See https://github.com/knative/pkg/issues/608

type KnativeRevision struct {
	Project           string
	Location          string
	ClusterName       string
	NamespaceName     string
	ServiceName       string
	ConfigurationName string
	RevisionName      string
}

func (kr *KnativeRevision) MonitoredResource() (resType string, labels map[string]string) {
	labels = map[string]string{
		metricskey.LabelProject:                  kr.Project,
		metricskey.LabelLocation:                 kr.Location,
		metricskey.LabelClusterName:              kr.ClusterName,
		metricskey.LabelNamespaceName:            kr.NamespaceName,
		metricskeyserving.LabelServiceName:       kr.ServiceName,
		metricskeyserving.LabelConfigurationName: kr.ConfigurationName,
		metricskeyserving.LabelRevisionName:      kr.RevisionName,
	}
	return "knative_revision", labels
}

func GetKnativeRevisionMonitoredResource(
	v *view.View, tags []tag.Tag, gm *monitoredresources.GcpMetadata) ([]tag.Tag, monitoredresource.Interface) {
	tagsMap := monitoredresources.GetTagsMap(tags)
	kr := &KnativeRevision{
		// The first three resource labels are from metadata.
		Project:     gm.Project,
		Location:    gm.Location,
		ClusterName: gm.Cluster,
		// The rest resource labels are from metrics labels.
		NamespaceName:     monitoredresources.ValueOrUnknown(metricskey.LabelNamespaceName, tagsMap),
		ServiceName:       monitoredresources.ValueOrUnknown(metricskeyserving.LabelServiceName, tagsMap),
		ConfigurationName: monitoredresources.ValueOrUnknown(metricskeyserving.LabelConfigurationName, tagsMap),
		RevisionName:      monitoredresources.ValueOrUnknown(metricskeyserving.LabelRevisionName, tagsMap),
	}

	var newTags []tag.Tag
	for _, t := range tags {
		// Keep the metrics labels that are not resource labels
		if !metricskeyserving.KnativeRevisionLabels.Has(t.Key.Name()) {
			newTags = append(newTags, t)
		}
	}

	return newTags, kr
}
