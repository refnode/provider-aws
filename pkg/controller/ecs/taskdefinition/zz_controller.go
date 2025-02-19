/*
Copyright 2021 The Crossplane Authors.

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

// Code generated by ack-generate. DO NOT EDIT.

package taskdefinition

import (
	"context"

	svcapi "github.com/aws/aws-sdk-go/service/ecs"
	svcsdk "github.com/aws/aws-sdk-go/service/ecs"
	svcsdkapi "github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane-contrib/provider-aws/apis/ecs/v1alpha1"
	awsclient "github.com/crossplane-contrib/provider-aws/pkg/clients"
)

const (
	errUnexpectedObject = "managed resource is not an TaskDefinition resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create TaskDefinition in AWS"
	errUpdate        = "cannot update TaskDefinition in AWS"
	errDescribe      = "failed to describe TaskDefinition"
	errDelete        = "failed to delete TaskDefinition"
)

type connector struct {
	kube client.Client
	opts []option
}

func (c *connector) Connect(ctx context.Context, mg cpresource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*svcapitypes.TaskDefinition)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	sess, err := awsclient.GetConfigV1(ctx, c.kube, mg, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, errors.Wrap(err, errCreateSession)
	}
	return newExternal(c.kube, svcapi.New(sess), c.opts), nil
}

func (e *external) Observe(ctx context.Context, mg cpresource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*svcapitypes.TaskDefinition)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateDescribeTaskDefinitionInput(cr)
	if err := e.preObserve(ctx, cr, input); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	resp, err := e.client.DescribeTaskDefinitionWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err := e.lateInitialize(&cr.Spec.ForProvider, resp); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "late-init failed")
	}
	GenerateTaskDefinition(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)

	upToDate, err := e.isUpToDate(cr, resp)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "isUpToDate check failed")
	}
	return e.postObserve(ctx, cr, resp, managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        upToDate,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil)
}

func (e *external) Create(ctx context.Context, mg cpresource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*svcapitypes.TaskDefinition)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Creating())
	input := GenerateRegisterTaskDefinitionInput(cr)
	if err := e.preCreate(ctx, cr, input); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	resp, err := e.client.RegisterTaskDefinitionWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, awsclient.Wrap(err, errCreate)
	}

	if resp.Tags != nil {
		f0 := []*svcapitypes.Tag{}
		for _, f0iter := range resp.Tags {
			f0elem := &svcapitypes.Tag{}
			if f0iter.Key != nil {
				f0elem.Key = f0iter.Key
			}
			if f0iter.Value != nil {
				f0elem.Value = f0iter.Value
			}
			f0 = append(f0, f0elem)
		}
		cr.Spec.ForProvider.Tags = f0
	} else {
		cr.Spec.ForProvider.Tags = nil
	}
	if resp.TaskDefinition != nil {
		f1 := &svcapitypes.TaskDefinition_SDK{}
		if resp.TaskDefinition.Compatibilities != nil {
			f1f0 := []*string{}
			for _, f1f0iter := range resp.TaskDefinition.Compatibilities {
				var f1f0elem string
				f1f0elem = *f1f0iter
				f1f0 = append(f1f0, &f1f0elem)
			}
			f1.Compatibilities = f1f0
		}
		if resp.TaskDefinition.ContainerDefinitions != nil {
			f1f1 := []*svcapitypes.ContainerDefinition{}
			for _, f1f1iter := range resp.TaskDefinition.ContainerDefinitions {
				f1f1elem := &svcapitypes.ContainerDefinition{}
				if f1f1iter.Command != nil {
					f1f1elemf0 := []*string{}
					for _, f1f1elemf0iter := range f1f1iter.Command {
						var f1f1elemf0elem string
						f1f1elemf0elem = *f1f1elemf0iter
						f1f1elemf0 = append(f1f1elemf0, &f1f1elemf0elem)
					}
					f1f1elem.Command = f1f1elemf0
				}
				if f1f1iter.Cpu != nil {
					f1f1elem.CPU = f1f1iter.Cpu
				}
				if f1f1iter.DependsOn != nil {
					f1f1elemf2 := []*svcapitypes.ContainerDependency{}
					for _, f1f1elemf2iter := range f1f1iter.DependsOn {
						f1f1elemf2elem := &svcapitypes.ContainerDependency{}
						if f1f1elemf2iter.Condition != nil {
							f1f1elemf2elem.Condition = f1f1elemf2iter.Condition
						}
						if f1f1elemf2iter.ContainerName != nil {
							f1f1elemf2elem.ContainerName = f1f1elemf2iter.ContainerName
						}
						f1f1elemf2 = append(f1f1elemf2, f1f1elemf2elem)
					}
					f1f1elem.DependsOn = f1f1elemf2
				}
				if f1f1iter.DisableNetworking != nil {
					f1f1elem.DisableNetworking = f1f1iter.DisableNetworking
				}
				if f1f1iter.DnsSearchDomains != nil {
					f1f1elemf4 := []*string{}
					for _, f1f1elemf4iter := range f1f1iter.DnsSearchDomains {
						var f1f1elemf4elem string
						f1f1elemf4elem = *f1f1elemf4iter
						f1f1elemf4 = append(f1f1elemf4, &f1f1elemf4elem)
					}
					f1f1elem.DNSSearchDomains = f1f1elemf4
				}
				if f1f1iter.DnsServers != nil {
					f1f1elemf5 := []*string{}
					for _, f1f1elemf5iter := range f1f1iter.DnsServers {
						var f1f1elemf5elem string
						f1f1elemf5elem = *f1f1elemf5iter
						f1f1elemf5 = append(f1f1elemf5, &f1f1elemf5elem)
					}
					f1f1elem.DNSServers = f1f1elemf5
				}
				if f1f1iter.DockerLabels != nil {
					f1f1elemf6 := map[string]*string{}
					for f1f1elemf6key, f1f1elemf6valiter := range f1f1iter.DockerLabels {
						var f1f1elemf6val string
						f1f1elemf6val = *f1f1elemf6valiter
						f1f1elemf6[f1f1elemf6key] = &f1f1elemf6val
					}
					f1f1elem.DockerLabels = f1f1elemf6
				}
				if f1f1iter.DockerSecurityOptions != nil {
					f1f1elemf7 := []*string{}
					for _, f1f1elemf7iter := range f1f1iter.DockerSecurityOptions {
						var f1f1elemf7elem string
						f1f1elemf7elem = *f1f1elemf7iter
						f1f1elemf7 = append(f1f1elemf7, &f1f1elemf7elem)
					}
					f1f1elem.DockerSecurityOptions = f1f1elemf7
				}
				if f1f1iter.EntryPoint != nil {
					f1f1elemf8 := []*string{}
					for _, f1f1elemf8iter := range f1f1iter.EntryPoint {
						var f1f1elemf8elem string
						f1f1elemf8elem = *f1f1elemf8iter
						f1f1elemf8 = append(f1f1elemf8, &f1f1elemf8elem)
					}
					f1f1elem.EntryPoint = f1f1elemf8
				}
				if f1f1iter.Environment != nil {
					f1f1elemf9 := []*svcapitypes.KeyValuePair{}
					for _, f1f1elemf9iter := range f1f1iter.Environment {
						f1f1elemf9elem := &svcapitypes.KeyValuePair{}
						if f1f1elemf9iter.Name != nil {
							f1f1elemf9elem.Name = f1f1elemf9iter.Name
						}
						if f1f1elemf9iter.Value != nil {
							f1f1elemf9elem.Value = f1f1elemf9iter.Value
						}
						f1f1elemf9 = append(f1f1elemf9, f1f1elemf9elem)
					}
					f1f1elem.Environment = f1f1elemf9
				}
				if f1f1iter.EnvironmentFiles != nil {
					f1f1elemf10 := []*svcapitypes.EnvironmentFile{}
					for _, f1f1elemf10iter := range f1f1iter.EnvironmentFiles {
						f1f1elemf10elem := &svcapitypes.EnvironmentFile{}
						if f1f1elemf10iter.Type != nil {
							f1f1elemf10elem.Type = f1f1elemf10iter.Type
						}
						if f1f1elemf10iter.Value != nil {
							f1f1elemf10elem.Value = f1f1elemf10iter.Value
						}
						f1f1elemf10 = append(f1f1elemf10, f1f1elemf10elem)
					}
					f1f1elem.EnvironmentFiles = f1f1elemf10
				}
				if f1f1iter.Essential != nil {
					f1f1elem.Essential = f1f1iter.Essential
				}
				if f1f1iter.ExtraHosts != nil {
					f1f1elemf12 := []*svcapitypes.HostEntry{}
					for _, f1f1elemf12iter := range f1f1iter.ExtraHosts {
						f1f1elemf12elem := &svcapitypes.HostEntry{}
						if f1f1elemf12iter.Hostname != nil {
							f1f1elemf12elem.Hostname = f1f1elemf12iter.Hostname
						}
						if f1f1elemf12iter.IpAddress != nil {
							f1f1elemf12elem.IPAddress = f1f1elemf12iter.IpAddress
						}
						f1f1elemf12 = append(f1f1elemf12, f1f1elemf12elem)
					}
					f1f1elem.ExtraHosts = f1f1elemf12
				}
				if f1f1iter.FirelensConfiguration != nil {
					f1f1elemf13 := &svcapitypes.FirelensConfiguration{}
					if f1f1iter.FirelensConfiguration.Options != nil {
						f1f1elemf13f0 := map[string]*string{}
						for f1f1elemf13f0key, f1f1elemf13f0valiter := range f1f1iter.FirelensConfiguration.Options {
							var f1f1elemf13f0val string
							f1f1elemf13f0val = *f1f1elemf13f0valiter
							f1f1elemf13f0[f1f1elemf13f0key] = &f1f1elemf13f0val
						}
						f1f1elemf13.Options = f1f1elemf13f0
					}
					if f1f1iter.FirelensConfiguration.Type != nil {
						f1f1elemf13.Type = f1f1iter.FirelensConfiguration.Type
					}
					f1f1elem.FirelensConfiguration = f1f1elemf13
				}
				if f1f1iter.HealthCheck != nil {
					f1f1elemf14 := &svcapitypes.HealthCheck{}
					if f1f1iter.HealthCheck.Command != nil {
						f1f1elemf14f0 := []*string{}
						for _, f1f1elemf14f0iter := range f1f1iter.HealthCheck.Command {
							var f1f1elemf14f0elem string
							f1f1elemf14f0elem = *f1f1elemf14f0iter
							f1f1elemf14f0 = append(f1f1elemf14f0, &f1f1elemf14f0elem)
						}
						f1f1elemf14.Command = f1f1elemf14f0
					}
					if f1f1iter.HealthCheck.Interval != nil {
						f1f1elemf14.Interval = f1f1iter.HealthCheck.Interval
					}
					if f1f1iter.HealthCheck.Retries != nil {
						f1f1elemf14.Retries = f1f1iter.HealthCheck.Retries
					}
					if f1f1iter.HealthCheck.StartPeriod != nil {
						f1f1elemf14.StartPeriod = f1f1iter.HealthCheck.StartPeriod
					}
					if f1f1iter.HealthCheck.Timeout != nil {
						f1f1elemf14.Timeout = f1f1iter.HealthCheck.Timeout
					}
					f1f1elem.HealthCheck = f1f1elemf14
				}
				if f1f1iter.Hostname != nil {
					f1f1elem.Hostname = f1f1iter.Hostname
				}
				if f1f1iter.Image != nil {
					f1f1elem.Image = f1f1iter.Image
				}
				if f1f1iter.Interactive != nil {
					f1f1elem.Interactive = f1f1iter.Interactive
				}
				if f1f1iter.Links != nil {
					f1f1elemf18 := []*string{}
					for _, f1f1elemf18iter := range f1f1iter.Links {
						var f1f1elemf18elem string
						f1f1elemf18elem = *f1f1elemf18iter
						f1f1elemf18 = append(f1f1elemf18, &f1f1elemf18elem)
					}
					f1f1elem.Links = f1f1elemf18
				}
				if f1f1iter.LinuxParameters != nil {
					f1f1elemf19 := &svcapitypes.LinuxParameters{}
					if f1f1iter.LinuxParameters.Capabilities != nil {
						f1f1elemf19f0 := &svcapitypes.KernelCapabilities{}
						if f1f1iter.LinuxParameters.Capabilities.Add != nil {
							f1f1elemf19f0f0 := []*string{}
							for _, f1f1elemf19f0f0iter := range f1f1iter.LinuxParameters.Capabilities.Add {
								var f1f1elemf19f0f0elem string
								f1f1elemf19f0f0elem = *f1f1elemf19f0f0iter
								f1f1elemf19f0f0 = append(f1f1elemf19f0f0, &f1f1elemf19f0f0elem)
							}
							f1f1elemf19f0.Add = f1f1elemf19f0f0
						}
						if f1f1iter.LinuxParameters.Capabilities.Drop != nil {
							f1f1elemf19f0f1 := []*string{}
							for _, f1f1elemf19f0f1iter := range f1f1iter.LinuxParameters.Capabilities.Drop {
								var f1f1elemf19f0f1elem string
								f1f1elemf19f0f1elem = *f1f1elemf19f0f1iter
								f1f1elemf19f0f1 = append(f1f1elemf19f0f1, &f1f1elemf19f0f1elem)
							}
							f1f1elemf19f0.Drop = f1f1elemf19f0f1
						}
						f1f1elemf19.Capabilities = f1f1elemf19f0
					}
					if f1f1iter.LinuxParameters.Devices != nil {
						f1f1elemf19f1 := []*svcapitypes.Device{}
						for _, f1f1elemf19f1iter := range f1f1iter.LinuxParameters.Devices {
							f1f1elemf19f1elem := &svcapitypes.Device{}
							if f1f1elemf19f1iter.ContainerPath != nil {
								f1f1elemf19f1elem.ContainerPath = f1f1elemf19f1iter.ContainerPath
							}
							if f1f1elemf19f1iter.HostPath != nil {
								f1f1elemf19f1elem.HostPath = f1f1elemf19f1iter.HostPath
							}
							if f1f1elemf19f1iter.Permissions != nil {
								f1f1elemf19f1elemf2 := []*string{}
								for _, f1f1elemf19f1elemf2iter := range f1f1elemf19f1iter.Permissions {
									var f1f1elemf19f1elemf2elem string
									f1f1elemf19f1elemf2elem = *f1f1elemf19f1elemf2iter
									f1f1elemf19f1elemf2 = append(f1f1elemf19f1elemf2, &f1f1elemf19f1elemf2elem)
								}
								f1f1elemf19f1elem.Permissions = f1f1elemf19f1elemf2
							}
							f1f1elemf19f1 = append(f1f1elemf19f1, f1f1elemf19f1elem)
						}
						f1f1elemf19.Devices = f1f1elemf19f1
					}
					if f1f1iter.LinuxParameters.InitProcessEnabled != nil {
						f1f1elemf19.InitProcessEnabled = f1f1iter.LinuxParameters.InitProcessEnabled
					}
					if f1f1iter.LinuxParameters.MaxSwap != nil {
						f1f1elemf19.MaxSwap = f1f1iter.LinuxParameters.MaxSwap
					}
					if f1f1iter.LinuxParameters.SharedMemorySize != nil {
						f1f1elemf19.SharedMemorySize = f1f1iter.LinuxParameters.SharedMemorySize
					}
					if f1f1iter.LinuxParameters.Swappiness != nil {
						f1f1elemf19.Swappiness = f1f1iter.LinuxParameters.Swappiness
					}
					if f1f1iter.LinuxParameters.Tmpfs != nil {
						f1f1elemf19f6 := []*svcapitypes.Tmpfs{}
						for _, f1f1elemf19f6iter := range f1f1iter.LinuxParameters.Tmpfs {
							f1f1elemf19f6elem := &svcapitypes.Tmpfs{}
							if f1f1elemf19f6iter.ContainerPath != nil {
								f1f1elemf19f6elem.ContainerPath = f1f1elemf19f6iter.ContainerPath
							}
							if f1f1elemf19f6iter.MountOptions != nil {
								f1f1elemf19f6elemf1 := []*string{}
								for _, f1f1elemf19f6elemf1iter := range f1f1elemf19f6iter.MountOptions {
									var f1f1elemf19f6elemf1elem string
									f1f1elemf19f6elemf1elem = *f1f1elemf19f6elemf1iter
									f1f1elemf19f6elemf1 = append(f1f1elemf19f6elemf1, &f1f1elemf19f6elemf1elem)
								}
								f1f1elemf19f6elem.MountOptions = f1f1elemf19f6elemf1
							}
							if f1f1elemf19f6iter.Size != nil {
								f1f1elemf19f6elem.Size = f1f1elemf19f6iter.Size
							}
							f1f1elemf19f6 = append(f1f1elemf19f6, f1f1elemf19f6elem)
						}
						f1f1elemf19.Tmpfs = f1f1elemf19f6
					}
					f1f1elem.LinuxParameters = f1f1elemf19
				}
				if f1f1iter.LogConfiguration != nil {
					f1f1elemf20 := &svcapitypes.LogConfiguration{}
					if f1f1iter.LogConfiguration.LogDriver != nil {
						f1f1elemf20.LogDriver = f1f1iter.LogConfiguration.LogDriver
					}
					if f1f1iter.LogConfiguration.Options != nil {
						f1f1elemf20f1 := map[string]*string{}
						for f1f1elemf20f1key, f1f1elemf20f1valiter := range f1f1iter.LogConfiguration.Options {
							var f1f1elemf20f1val string
							f1f1elemf20f1val = *f1f1elemf20f1valiter
							f1f1elemf20f1[f1f1elemf20f1key] = &f1f1elemf20f1val
						}
						f1f1elemf20.Options = f1f1elemf20f1
					}
					if f1f1iter.LogConfiguration.SecretOptions != nil {
						f1f1elemf20f2 := []*svcapitypes.Secret{}
						for _, f1f1elemf20f2iter := range f1f1iter.LogConfiguration.SecretOptions {
							f1f1elemf20f2elem := &svcapitypes.Secret{}
							if f1f1elemf20f2iter.Name != nil {
								f1f1elemf20f2elem.Name = f1f1elemf20f2iter.Name
							}
							if f1f1elemf20f2iter.ValueFrom != nil {
								f1f1elemf20f2elem.ValueFrom = f1f1elemf20f2iter.ValueFrom
							}
							f1f1elemf20f2 = append(f1f1elemf20f2, f1f1elemf20f2elem)
						}
						f1f1elemf20.SecretOptions = f1f1elemf20f2
					}
					f1f1elem.LogConfiguration = f1f1elemf20
				}
				if f1f1iter.Memory != nil {
					f1f1elem.Memory = f1f1iter.Memory
				}
				if f1f1iter.MemoryReservation != nil {
					f1f1elem.MemoryReservation = f1f1iter.MemoryReservation
				}
				if f1f1iter.MountPoints != nil {
					f1f1elemf23 := []*svcapitypes.MountPoint{}
					for _, f1f1elemf23iter := range f1f1iter.MountPoints {
						f1f1elemf23elem := &svcapitypes.MountPoint{}
						if f1f1elemf23iter.ContainerPath != nil {
							f1f1elemf23elem.ContainerPath = f1f1elemf23iter.ContainerPath
						}
						if f1f1elemf23iter.ReadOnly != nil {
							f1f1elemf23elem.ReadOnly = f1f1elemf23iter.ReadOnly
						}
						if f1f1elemf23iter.SourceVolume != nil {
							f1f1elemf23elem.SourceVolume = f1f1elemf23iter.SourceVolume
						}
						f1f1elemf23 = append(f1f1elemf23, f1f1elemf23elem)
					}
					f1f1elem.MountPoints = f1f1elemf23
				}
				if f1f1iter.Name != nil {
					f1f1elem.Name = f1f1iter.Name
				}
				if f1f1iter.PortMappings != nil {
					f1f1elemf25 := []*svcapitypes.PortMapping{}
					for _, f1f1elemf25iter := range f1f1iter.PortMappings {
						f1f1elemf25elem := &svcapitypes.PortMapping{}
						if f1f1elemf25iter.AppProtocol != nil {
							f1f1elemf25elem.AppProtocol = f1f1elemf25iter.AppProtocol
						}
						if f1f1elemf25iter.ContainerPort != nil {
							f1f1elemf25elem.ContainerPort = f1f1elemf25iter.ContainerPort
						}
						if f1f1elemf25iter.ContainerPortRange != nil {
							f1f1elemf25elem.ContainerPortRange = f1f1elemf25iter.ContainerPortRange
						}
						if f1f1elemf25iter.HostPort != nil {
							f1f1elemf25elem.HostPort = f1f1elemf25iter.HostPort
						}
						if f1f1elemf25iter.Name != nil {
							f1f1elemf25elem.Name = f1f1elemf25iter.Name
						}
						if f1f1elemf25iter.Protocol != nil {
							f1f1elemf25elem.Protocol = f1f1elemf25iter.Protocol
						}
						f1f1elemf25 = append(f1f1elemf25, f1f1elemf25elem)
					}
					f1f1elem.PortMappings = f1f1elemf25
				}
				if f1f1iter.Privileged != nil {
					f1f1elem.Privileged = f1f1iter.Privileged
				}
				if f1f1iter.PseudoTerminal != nil {
					f1f1elem.PseudoTerminal = f1f1iter.PseudoTerminal
				}
				if f1f1iter.ReadonlyRootFilesystem != nil {
					f1f1elem.ReadonlyRootFilesystem = f1f1iter.ReadonlyRootFilesystem
				}
				if f1f1iter.RepositoryCredentials != nil {
					f1f1elemf29 := &svcapitypes.RepositoryCredentials{}
					if f1f1iter.RepositoryCredentials.CredentialsParameter != nil {
						f1f1elemf29.CredentialsParameter = f1f1iter.RepositoryCredentials.CredentialsParameter
					}
					f1f1elem.RepositoryCredentials = f1f1elemf29
				}
				if f1f1iter.ResourceRequirements != nil {
					f1f1elemf30 := []*svcapitypes.ResourceRequirement{}
					for _, f1f1elemf30iter := range f1f1iter.ResourceRequirements {
						f1f1elemf30elem := &svcapitypes.ResourceRequirement{}
						if f1f1elemf30iter.Type != nil {
							f1f1elemf30elem.Type = f1f1elemf30iter.Type
						}
						if f1f1elemf30iter.Value != nil {
							f1f1elemf30elem.Value = f1f1elemf30iter.Value
						}
						f1f1elemf30 = append(f1f1elemf30, f1f1elemf30elem)
					}
					f1f1elem.ResourceRequirements = f1f1elemf30
				}
				if f1f1iter.Secrets != nil {
					f1f1elemf31 := []*svcapitypes.Secret{}
					for _, f1f1elemf31iter := range f1f1iter.Secrets {
						f1f1elemf31elem := &svcapitypes.Secret{}
						if f1f1elemf31iter.Name != nil {
							f1f1elemf31elem.Name = f1f1elemf31iter.Name
						}
						if f1f1elemf31iter.ValueFrom != nil {
							f1f1elemf31elem.ValueFrom = f1f1elemf31iter.ValueFrom
						}
						f1f1elemf31 = append(f1f1elemf31, f1f1elemf31elem)
					}
					f1f1elem.Secrets = f1f1elemf31
				}
				if f1f1iter.StartTimeout != nil {
					f1f1elem.StartTimeout = f1f1iter.StartTimeout
				}
				if f1f1iter.StopTimeout != nil {
					f1f1elem.StopTimeout = f1f1iter.StopTimeout
				}
				if f1f1iter.SystemControls != nil {
					f1f1elemf34 := []*svcapitypes.SystemControl{}
					for _, f1f1elemf34iter := range f1f1iter.SystemControls {
						f1f1elemf34elem := &svcapitypes.SystemControl{}
						if f1f1elemf34iter.Namespace != nil {
							f1f1elemf34elem.Namespace = f1f1elemf34iter.Namespace
						}
						if f1f1elemf34iter.Value != nil {
							f1f1elemf34elem.Value = f1f1elemf34iter.Value
						}
						f1f1elemf34 = append(f1f1elemf34, f1f1elemf34elem)
					}
					f1f1elem.SystemControls = f1f1elemf34
				}
				if f1f1iter.Ulimits != nil {
					f1f1elemf35 := []*svcapitypes.Ulimit{}
					for _, f1f1elemf35iter := range f1f1iter.Ulimits {
						f1f1elemf35elem := &svcapitypes.Ulimit{}
						if f1f1elemf35iter.HardLimit != nil {
							f1f1elemf35elem.HardLimit = f1f1elemf35iter.HardLimit
						}
						if f1f1elemf35iter.Name != nil {
							f1f1elemf35elem.Name = f1f1elemf35iter.Name
						}
						if f1f1elemf35iter.SoftLimit != nil {
							f1f1elemf35elem.SoftLimit = f1f1elemf35iter.SoftLimit
						}
						f1f1elemf35 = append(f1f1elemf35, f1f1elemf35elem)
					}
					f1f1elem.Ulimits = f1f1elemf35
				}
				if f1f1iter.User != nil {
					f1f1elem.User = f1f1iter.User
				}
				if f1f1iter.VolumesFrom != nil {
					f1f1elemf37 := []*svcapitypes.VolumeFrom{}
					for _, f1f1elemf37iter := range f1f1iter.VolumesFrom {
						f1f1elemf37elem := &svcapitypes.VolumeFrom{}
						if f1f1elemf37iter.ReadOnly != nil {
							f1f1elemf37elem.ReadOnly = f1f1elemf37iter.ReadOnly
						}
						if f1f1elemf37iter.SourceContainer != nil {
							f1f1elemf37elem.SourceContainer = f1f1elemf37iter.SourceContainer
						}
						f1f1elemf37 = append(f1f1elemf37, f1f1elemf37elem)
					}
					f1f1elem.VolumesFrom = f1f1elemf37
				}
				if f1f1iter.WorkingDirectory != nil {
					f1f1elem.WorkingDirectory = f1f1iter.WorkingDirectory
				}
				f1f1 = append(f1f1, f1f1elem)
			}
			f1.ContainerDefinitions = f1f1
		}
		if resp.TaskDefinition.Cpu != nil {
			f1.CPU = resp.TaskDefinition.Cpu
		}
		if resp.TaskDefinition.DeregisteredAt != nil {
			f1.DeregisteredAt = &metav1.Time{*resp.TaskDefinition.DeregisteredAt}
		}
		if resp.TaskDefinition.EphemeralStorage != nil {
			f1f4 := &svcapitypes.EphemeralStorage{}
			if resp.TaskDefinition.EphemeralStorage.SizeInGiB != nil {
				f1f4.SizeInGiB = resp.TaskDefinition.EphemeralStorage.SizeInGiB
			}
			f1.EphemeralStorage = f1f4
		}
		if resp.TaskDefinition.ExecutionRoleArn != nil {
			f1.ExecutionRoleARN = resp.TaskDefinition.ExecutionRoleArn
		}
		if resp.TaskDefinition.Family != nil {
			f1.Family = resp.TaskDefinition.Family
		}
		if resp.TaskDefinition.InferenceAccelerators != nil {
			f1f7 := []*svcapitypes.InferenceAccelerator{}
			for _, f1f7iter := range resp.TaskDefinition.InferenceAccelerators {
				f1f7elem := &svcapitypes.InferenceAccelerator{}
				if f1f7iter.DeviceName != nil {
					f1f7elem.DeviceName = f1f7iter.DeviceName
				}
				if f1f7iter.DeviceType != nil {
					f1f7elem.DeviceType = f1f7iter.DeviceType
				}
				f1f7 = append(f1f7, f1f7elem)
			}
			f1.InferenceAccelerators = f1f7
		}
		if resp.TaskDefinition.IpcMode != nil {
			f1.IPCMode = resp.TaskDefinition.IpcMode
		}
		if resp.TaskDefinition.Memory != nil {
			f1.Memory = resp.TaskDefinition.Memory
		}
		if resp.TaskDefinition.NetworkMode != nil {
			f1.NetworkMode = resp.TaskDefinition.NetworkMode
		}
		if resp.TaskDefinition.PidMode != nil {
			f1.PIDMode = resp.TaskDefinition.PidMode
		}
		if resp.TaskDefinition.PlacementConstraints != nil {
			f1f12 := []*svcapitypes.TaskDefinitionPlacementConstraint{}
			for _, f1f12iter := range resp.TaskDefinition.PlacementConstraints {
				f1f12elem := &svcapitypes.TaskDefinitionPlacementConstraint{}
				if f1f12iter.Expression != nil {
					f1f12elem.Expression = f1f12iter.Expression
				}
				if f1f12iter.Type != nil {
					f1f12elem.Type = f1f12iter.Type
				}
				f1f12 = append(f1f12, f1f12elem)
			}
			f1.PlacementConstraints = f1f12
		}
		if resp.TaskDefinition.ProxyConfiguration != nil {
			f1f13 := &svcapitypes.ProxyConfiguration{}
			if resp.TaskDefinition.ProxyConfiguration.ContainerName != nil {
				f1f13.ContainerName = resp.TaskDefinition.ProxyConfiguration.ContainerName
			}
			if resp.TaskDefinition.ProxyConfiguration.Properties != nil {
				f1f13f1 := []*svcapitypes.KeyValuePair{}
				for _, f1f13f1iter := range resp.TaskDefinition.ProxyConfiguration.Properties {
					f1f13f1elem := &svcapitypes.KeyValuePair{}
					if f1f13f1iter.Name != nil {
						f1f13f1elem.Name = f1f13f1iter.Name
					}
					if f1f13f1iter.Value != nil {
						f1f13f1elem.Value = f1f13f1iter.Value
					}
					f1f13f1 = append(f1f13f1, f1f13f1elem)
				}
				f1f13.Properties = f1f13f1
			}
			if resp.TaskDefinition.ProxyConfiguration.Type != nil {
				f1f13.Type = resp.TaskDefinition.ProxyConfiguration.Type
			}
			f1.ProxyConfiguration = f1f13
		}
		if resp.TaskDefinition.RegisteredAt != nil {
			f1.RegisteredAt = &metav1.Time{*resp.TaskDefinition.RegisteredAt}
		}
		if resp.TaskDefinition.RegisteredBy != nil {
			f1.RegisteredBy = resp.TaskDefinition.RegisteredBy
		}
		if resp.TaskDefinition.RequiresAttributes != nil {
			f1f16 := []*svcapitypes.Attribute{}
			for _, f1f16iter := range resp.TaskDefinition.RequiresAttributes {
				f1f16elem := &svcapitypes.Attribute{}
				if f1f16iter.Name != nil {
					f1f16elem.Name = f1f16iter.Name
				}
				if f1f16iter.TargetId != nil {
					f1f16elem.TargetID = f1f16iter.TargetId
				}
				if f1f16iter.TargetType != nil {
					f1f16elem.TargetType = f1f16iter.TargetType
				}
				if f1f16iter.Value != nil {
					f1f16elem.Value = f1f16iter.Value
				}
				f1f16 = append(f1f16, f1f16elem)
			}
			f1.RequiresAttributes = f1f16
		}
		if resp.TaskDefinition.RequiresCompatibilities != nil {
			f1f17 := []*string{}
			for _, f1f17iter := range resp.TaskDefinition.RequiresCompatibilities {
				var f1f17elem string
				f1f17elem = *f1f17iter
				f1f17 = append(f1f17, &f1f17elem)
			}
			f1.RequiresCompatibilities = f1f17
		}
		if resp.TaskDefinition.Revision != nil {
			f1.Revision = resp.TaskDefinition.Revision
		}
		if resp.TaskDefinition.RuntimePlatform != nil {
			f1f19 := &svcapitypes.RuntimePlatform{}
			if resp.TaskDefinition.RuntimePlatform.CpuArchitecture != nil {
				f1f19.CPUArchitecture = resp.TaskDefinition.RuntimePlatform.CpuArchitecture
			}
			if resp.TaskDefinition.RuntimePlatform.OperatingSystemFamily != nil {
				f1f19.OperatingSystemFamily = resp.TaskDefinition.RuntimePlatform.OperatingSystemFamily
			}
			f1.RuntimePlatform = f1f19
		}
		if resp.TaskDefinition.Status != nil {
			f1.Status = resp.TaskDefinition.Status
		}
		if resp.TaskDefinition.TaskDefinitionArn != nil {
			f1.TaskDefinitionARN = resp.TaskDefinition.TaskDefinitionArn
		}
		if resp.TaskDefinition.TaskRoleArn != nil {
			f1.TaskRoleARN = resp.TaskDefinition.TaskRoleArn
		}
		if resp.TaskDefinition.Volumes != nil {
			f1f23 := []*svcapitypes.Volume{}
			for _, f1f23iter := range resp.TaskDefinition.Volumes {
				f1f23elem := &svcapitypes.Volume{}
				if f1f23iter.DockerVolumeConfiguration != nil {
					f1f23elemf0 := &svcapitypes.DockerVolumeConfiguration{}
					if f1f23iter.DockerVolumeConfiguration.Autoprovision != nil {
						f1f23elemf0.Autoprovision = f1f23iter.DockerVolumeConfiguration.Autoprovision
					}
					if f1f23iter.DockerVolumeConfiguration.Driver != nil {
						f1f23elemf0.Driver = f1f23iter.DockerVolumeConfiguration.Driver
					}
					if f1f23iter.DockerVolumeConfiguration.DriverOpts != nil {
						f1f23elemf0f2 := map[string]*string{}
						for f1f23elemf0f2key, f1f23elemf0f2valiter := range f1f23iter.DockerVolumeConfiguration.DriverOpts {
							var f1f23elemf0f2val string
							f1f23elemf0f2val = *f1f23elemf0f2valiter
							f1f23elemf0f2[f1f23elemf0f2key] = &f1f23elemf0f2val
						}
						f1f23elemf0.DriverOpts = f1f23elemf0f2
					}
					if f1f23iter.DockerVolumeConfiguration.Labels != nil {
						f1f23elemf0f3 := map[string]*string{}
						for f1f23elemf0f3key, f1f23elemf0f3valiter := range f1f23iter.DockerVolumeConfiguration.Labels {
							var f1f23elemf0f3val string
							f1f23elemf0f3val = *f1f23elemf0f3valiter
							f1f23elemf0f3[f1f23elemf0f3key] = &f1f23elemf0f3val
						}
						f1f23elemf0.Labels = f1f23elemf0f3
					}
					if f1f23iter.DockerVolumeConfiguration.Scope != nil {
						f1f23elemf0.Scope = f1f23iter.DockerVolumeConfiguration.Scope
					}
					f1f23elem.DockerVolumeConfiguration = f1f23elemf0
				}
				if f1f23iter.EfsVolumeConfiguration != nil {
					f1f23elemf1 := &svcapitypes.EFSVolumeConfiguration{}
					if f1f23iter.EfsVolumeConfiguration.AuthorizationConfig != nil {
						f1f23elemf1f0 := &svcapitypes.EFSAuthorizationConfig{}
						if f1f23iter.EfsVolumeConfiguration.AuthorizationConfig.AccessPointId != nil {
							f1f23elemf1f0.AccessPointID = f1f23iter.EfsVolumeConfiguration.AuthorizationConfig.AccessPointId
						}
						if f1f23iter.EfsVolumeConfiguration.AuthorizationConfig.Iam != nil {
							f1f23elemf1f0.IAM = f1f23iter.EfsVolumeConfiguration.AuthorizationConfig.Iam
						}
						f1f23elemf1.AuthorizationConfig = f1f23elemf1f0
					}
					if f1f23iter.EfsVolumeConfiguration.FileSystemId != nil {
						f1f23elemf1.FileSystemID = f1f23iter.EfsVolumeConfiguration.FileSystemId
					}
					if f1f23iter.EfsVolumeConfiguration.RootDirectory != nil {
						f1f23elemf1.RootDirectory = f1f23iter.EfsVolumeConfiguration.RootDirectory
					}
					if f1f23iter.EfsVolumeConfiguration.TransitEncryption != nil {
						f1f23elemf1.TransitEncryption = f1f23iter.EfsVolumeConfiguration.TransitEncryption
					}
					if f1f23iter.EfsVolumeConfiguration.TransitEncryptionPort != nil {
						f1f23elemf1.TransitEncryptionPort = f1f23iter.EfsVolumeConfiguration.TransitEncryptionPort
					}
					f1f23elem.EFSVolumeConfiguration = f1f23elemf1
				}
				if f1f23iter.FsxWindowsFileServerVolumeConfiguration != nil {
					f1f23elemf2 := &svcapitypes.FSxWindowsFileServerVolumeConfiguration{}
					if f1f23iter.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig != nil {
						f1f23elemf2f0 := &svcapitypes.FSxWindowsFileServerAuthorizationConfig{}
						if f1f23iter.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig.CredentialsParameter != nil {
							f1f23elemf2f0.CredentialsParameter = f1f23iter.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig.CredentialsParameter
						}
						if f1f23iter.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig.Domain != nil {
							f1f23elemf2f0.Domain = f1f23iter.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig.Domain
						}
						f1f23elemf2.AuthorizationConfig = f1f23elemf2f0
					}
					if f1f23iter.FsxWindowsFileServerVolumeConfiguration.FileSystemId != nil {
						f1f23elemf2.FileSystemID = f1f23iter.FsxWindowsFileServerVolumeConfiguration.FileSystemId
					}
					if f1f23iter.FsxWindowsFileServerVolumeConfiguration.RootDirectory != nil {
						f1f23elemf2.RootDirectory = f1f23iter.FsxWindowsFileServerVolumeConfiguration.RootDirectory
					}
					f1f23elem.FsxWindowsFileServerVolumeConfiguration = f1f23elemf2
				}
				if f1f23iter.Host != nil {
					f1f23elemf3 := &svcapitypes.HostVolumeProperties{}
					if f1f23iter.Host.SourcePath != nil {
						f1f23elemf3.SourcePath = f1f23iter.Host.SourcePath
					}
					f1f23elem.Host = f1f23elemf3
				}
				if f1f23iter.Name != nil {
					f1f23elem.Name = f1f23iter.Name
				}
				f1f23 = append(f1f23, f1f23elem)
			}
			f1.Volumes = f1f23
		}
		cr.Status.AtProvider.TaskDefinition = f1
	} else {
		cr.Status.AtProvider.TaskDefinition = nil
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
	return e.update(ctx, mg)

}

func (e *external) Delete(ctx context.Context, mg cpresource.Managed) error {
	cr, ok := mg.(*svcapitypes.TaskDefinition)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeregisterTaskDefinitionInput(cr)
	ignore, err := e.preDelete(ctx, cr, input)
	if err != nil {
		return errors.Wrap(err, "pre-delete failed")
	}
	if ignore {
		return nil
	}
	resp, err := e.client.DeregisterTaskDefinitionWithContext(ctx, input)
	return e.postDelete(ctx, cr, resp, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDelete))
}

type option func(*external)

func newExternal(kube client.Client, client svcsdkapi.ECSAPI, opts []option) *external {
	e := &external{
		kube:           kube,
		client:         client,
		preObserve:     nopPreObserve,
		postObserve:    nopPostObserve,
		lateInitialize: nopLateInitialize,
		isUpToDate:     alwaysUpToDate,
		preCreate:      nopPreCreate,
		postCreate:     nopPostCreate,
		preDelete:      nopPreDelete,
		postDelete:     nopPostDelete,
		update:         nopUpdate,
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

type external struct {
	kube           client.Client
	client         svcsdkapi.ECSAPI
	preObserve     func(context.Context, *svcapitypes.TaskDefinition, *svcsdk.DescribeTaskDefinitionInput) error
	postObserve    func(context.Context, *svcapitypes.TaskDefinition, *svcsdk.DescribeTaskDefinitionOutput, managed.ExternalObservation, error) (managed.ExternalObservation, error)
	lateInitialize func(*svcapitypes.TaskDefinitionParameters, *svcsdk.DescribeTaskDefinitionOutput) error
	isUpToDate     func(*svcapitypes.TaskDefinition, *svcsdk.DescribeTaskDefinitionOutput) (bool, error)
	preCreate      func(context.Context, *svcapitypes.TaskDefinition, *svcsdk.RegisterTaskDefinitionInput) error
	postCreate     func(context.Context, *svcapitypes.TaskDefinition, *svcsdk.RegisterTaskDefinitionOutput, managed.ExternalCreation, error) (managed.ExternalCreation, error)
	preDelete      func(context.Context, *svcapitypes.TaskDefinition, *svcsdk.DeregisterTaskDefinitionInput) (bool, error)
	postDelete     func(context.Context, *svcapitypes.TaskDefinition, *svcsdk.DeregisterTaskDefinitionOutput, error) error
	update         func(context.Context, cpresource.Managed) (managed.ExternalUpdate, error)
}

func nopPreObserve(context.Context, *svcapitypes.TaskDefinition, *svcsdk.DescribeTaskDefinitionInput) error {
	return nil
}

func nopPostObserve(_ context.Context, _ *svcapitypes.TaskDefinition, _ *svcsdk.DescribeTaskDefinitionOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	return obs, err
}
func nopLateInitialize(*svcapitypes.TaskDefinitionParameters, *svcsdk.DescribeTaskDefinitionOutput) error {
	return nil
}
func alwaysUpToDate(*svcapitypes.TaskDefinition, *svcsdk.DescribeTaskDefinitionOutput) (bool, error) {
	return true, nil
}

func nopPreCreate(context.Context, *svcapitypes.TaskDefinition, *svcsdk.RegisterTaskDefinitionInput) error {
	return nil
}
func nopPostCreate(_ context.Context, _ *svcapitypes.TaskDefinition, _ *svcsdk.RegisterTaskDefinitionOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	return cre, err
}
func nopPreDelete(context.Context, *svcapitypes.TaskDefinition, *svcsdk.DeregisterTaskDefinitionInput) (bool, error) {
	return false, nil
}
func nopPostDelete(_ context.Context, _ *svcapitypes.TaskDefinition, _ *svcsdk.DeregisterTaskDefinitionOutput, err error) error {
	return err
}
func nopUpdate(context.Context, cpresource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}
