/*
Copyright 2020 The Tekton Authors

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	resourcev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
	scheme "github.com/tektoncd/pipeline/pkg/client/resource/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// PipelineResourcesGetter has a method to return a PipelineResourceInterface.
// A group's client should implement this interface.
type PipelineResourcesGetter interface {
	PipelineResources(namespace string) PipelineResourceInterface
}

// PipelineResourceInterface has methods to work with PipelineResource resources.
type PipelineResourceInterface interface {
	Create(ctx context.Context, pipelineResource *resourcev1alpha1.PipelineResource, opts v1.CreateOptions) (*resourcev1alpha1.PipelineResource, error)
	Update(ctx context.Context, pipelineResource *resourcev1alpha1.PipelineResource, opts v1.UpdateOptions) (*resourcev1alpha1.PipelineResource, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*resourcev1alpha1.PipelineResource, error)
	List(ctx context.Context, opts v1.ListOptions) (*resourcev1alpha1.PipelineResourceList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *resourcev1alpha1.PipelineResource, err error)
	PipelineResourceExpansion
}

// pipelineResources implements PipelineResourceInterface
type pipelineResources struct {
	*gentype.ClientWithList[*resourcev1alpha1.PipelineResource, *resourcev1alpha1.PipelineResourceList]
}

// newPipelineResources returns a PipelineResources
func newPipelineResources(c *TektonV1alpha1Client, namespace string) *pipelineResources {
	return &pipelineResources{
		gentype.NewClientWithList[*resourcev1alpha1.PipelineResource, *resourcev1alpha1.PipelineResourceList](
			"pipelineresources",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *resourcev1alpha1.PipelineResource { return &resourcev1alpha1.PipelineResource{} },
			func() *resourcev1alpha1.PipelineResourceList { return &resourcev1alpha1.PipelineResourceList{} },
		),
	}
}
