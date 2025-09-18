/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package testing

import (
	"context"
	"net"

	ffv1 "github.com/innabox/fulfillment-common/api/fulfillment/v1"
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
)

// Server is a gRPC server used only for tests.
type Server struct {
	listener net.Listener
	server   *grpc.Server
}

// NewServer creates a new gRPC server that listens in a randomly selected port in the local host.
func NewServer() *Server {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	Expect(err).ToNot(HaveOccurred())
	server := grpc.NewServer()
	return &Server{
		listener: listener,
		server:   server,
	}
}

// Adress returns the address where the server is listening.
func (s *Server) Address() string {
	return s.listener.Addr().String()
}

// Registrar returns the registrar that can be used to register server implementations.
func (s *Server) Registrar() grpc.ServiceRegistrar {
	return s.server
}

// Start starts the server. This needs to be done after registering all server implementations, and before trying to
// call any of them.
func (s *Server) Start() {
	go func() {
		defer GinkgoRecover()
		err := s.server.Serve(s.listener)
		Expect(err).ToNot(HaveOccurred())
	}()
}

// Stop stops the server, closing all connections and releasing all the resources it was using.
func (s *Server) Stop() {
	s.server.Stop()
}

// Make sure that we implement the interface.
var _ ffv1.ClustersServer = (*ClustersServerFuncs)(nil)

// ClustersServerFuncs is an implementation of the clusters server that uses configurable functions to implement the
// methods.
type ClustersServerFuncs struct {
	ffv1.UnimplementedClustersServer

	CreateFunc               func(context.Context, *ffv1.ClustersCreateRequest) (*ffv1.ClustersCreateResponse, error)
	DeleteFunc               func(context.Context, *ffv1.ClustersDeleteRequest) (*ffv1.ClustersDeleteResponse, error)
	GetFunc                  func(context.Context, *ffv1.ClustersGetRequest) (*ffv1.ClustersGetResponse, error)
	ListFunc                 func(context.Context, *ffv1.ClustersListRequest) (*ffv1.ClustersListResponse, error)
	GetKubeconfigFunc        func(context.Context, *ffv1.ClustersGetKubeconfigRequest) (*ffv1.ClustersGetKubeconfigResponse, error)
	GetKubeconfigViaHttpFunc func(context.Context, *ffv1.ClustersGetKubeconfigViaHttpRequest) (*httpbody.HttpBody, error)
	UpdateFunc               func(context.Context, *ffv1.ClustersUpdateRequest) (*ffv1.ClustersUpdateResponse, error)
}

func (s *ClustersServerFuncs) Create(ctx context.Context,
	request *ffv1.ClustersCreateRequest) (response *ffv1.ClustersCreateResponse, err error) {
	response, err = s.CreateFunc(ctx, request)
	return
}

func (s *ClustersServerFuncs) Delete(ctx context.Context,
	request *ffv1.ClustersDeleteRequest) (response *ffv1.ClustersDeleteResponse, err error) {
	response, err = s.DeleteFunc(ctx, request)
	return
}

func (s *ClustersServerFuncs) Get(ctx context.Context,
	request *ffv1.ClustersGetRequest) (response *ffv1.ClustersGetResponse, err error) {
	response, err = s.GetFunc(ctx, request)
	return
}

func (s *ClustersServerFuncs) GetKubeconfig(ctx context.Context,
	request *ffv1.ClustersGetKubeconfigRequest) (response *ffv1.ClustersGetKubeconfigResponse, err error) {
	response, err = s.GetKubeconfigFunc(ctx, request)
	return
}

func (s *ClustersServerFuncs) GetKubeconfigViaHttp(ctx context.Context,
	request *ffv1.ClustersGetKubeconfigViaHttpRequest) (response *httpbody.HttpBody, err error) {
	response, err = s.GetKubeconfigViaHttpFunc(ctx, request)
	return
}

func (s *ClustersServerFuncs) List(ctx context.Context,
	request *ffv1.ClustersListRequest) (response *ffv1.ClustersListResponse, err error) {
	response, err = s.ListFunc(ctx, request)
	return
}

func (s *ClustersServerFuncs) Update(ctx context.Context,
	request *ffv1.ClustersUpdateRequest) (response *ffv1.ClustersUpdateResponse, err error) {
	response, err = s.UpdateFunc(ctx, request)
	return
}

// Make sure that we implement the interface.
var _ ffv1.HostsServer = (*HostsServerFuncs)(nil)

// HostsServerFuncs is an implementation of the hosts server that uses configurable functions to implement the
// methods.
type HostsServerFuncs struct {
	ffv1.UnimplementedHostsServer

	CreateFunc func(context.Context, *ffv1.HostsCreateRequest) (*ffv1.HostsCreateResponse, error)
	DeleteFunc func(context.Context, *ffv1.HostsDeleteRequest) (*ffv1.HostsDeleteResponse, error)
	GetFunc    func(context.Context, *ffv1.HostsGetRequest) (*ffv1.HostsGetResponse, error)
	ListFunc   func(context.Context, *ffv1.HostsListRequest) (*ffv1.HostsListResponse, error)
	UpdateFunc func(context.Context, *ffv1.HostsUpdateRequest) (*ffv1.HostsUpdateResponse, error)
}

func (s *HostsServerFuncs) Create(ctx context.Context,
	request *ffv1.HostsCreateRequest) (response *ffv1.HostsCreateResponse, err error) {
	response, err = s.CreateFunc(ctx, request)
	return
}

func (s *HostsServerFuncs) Delete(ctx context.Context,
	request *ffv1.HostsDeleteRequest) (response *ffv1.HostsDeleteResponse, err error) {
	response, err = s.DeleteFunc(ctx, request)
	return
}

func (s *HostsServerFuncs) Get(ctx context.Context,
	request *ffv1.HostsGetRequest) (response *ffv1.HostsGetResponse, err error) {
	response, err = s.GetFunc(ctx, request)
	return
}

func (s *HostsServerFuncs) List(ctx context.Context,
	request *ffv1.HostsListRequest) (response *ffv1.HostsListResponse, err error) {
	response, err = s.ListFunc(ctx, request)
	return
}

func (s *HostsServerFuncs) Update(ctx context.Context,
	request *ffv1.HostsUpdateRequest) (response *ffv1.HostsUpdateResponse, err error) {
	response, err = s.UpdateFunc(ctx, request)
	return
}

// Make sure that we implement the interface.
var _ ffv1.HostPoolsServer = (*HostPoolsServerFuncs)(nil)

// HostPoolsServerFuncs is an implementation of the host pools server that uses configurable functions to implement the
// methods.
type HostPoolsServerFuncs struct {
	ffv1.UnimplementedHostPoolsServer

	CreateFunc func(context.Context, *ffv1.HostPoolsCreateRequest) (*ffv1.HostPoolsCreateResponse, error)
	DeleteFunc func(context.Context, *ffv1.HostPoolsDeleteRequest) (*ffv1.HostPoolsDeleteResponse, error)
	GetFunc    func(context.Context, *ffv1.HostPoolsGetRequest) (*ffv1.HostPoolsGetResponse, error)
	ListFunc   func(context.Context, *ffv1.HostPoolsListRequest) (*ffv1.HostPoolsListResponse, error)
	UpdateFunc func(context.Context, *ffv1.HostPoolsUpdateRequest) (*ffv1.HostPoolsUpdateResponse, error)
}

func (s *HostPoolsServerFuncs) Create(ctx context.Context,
	request *ffv1.HostPoolsCreateRequest) (response *ffv1.HostPoolsCreateResponse, err error) {
	response, err = s.CreateFunc(ctx, request)
	return
}

func (s *HostPoolsServerFuncs) Delete(ctx context.Context,
	request *ffv1.HostPoolsDeleteRequest) (response *ffv1.HostPoolsDeleteResponse, err error) {
	response, err = s.DeleteFunc(ctx, request)
	return
}

func (s *HostPoolsServerFuncs) Get(ctx context.Context,
	request *ffv1.HostPoolsGetRequest) (response *ffv1.HostPoolsGetResponse, err error) {
	response, err = s.GetFunc(ctx, request)
	return
}

func (s *HostPoolsServerFuncs) List(ctx context.Context,
	request *ffv1.HostPoolsListRequest) (response *ffv1.HostPoolsListResponse, err error) {
	response, err = s.ListFunc(ctx, request)
	return
}

func (s *HostPoolsServerFuncs) Update(ctx context.Context,
	request *ffv1.HostPoolsUpdateRequest) (response *ffv1.HostPoolsUpdateResponse, err error) {
	response, err = s.UpdateFunc(ctx, request)
	return
}
