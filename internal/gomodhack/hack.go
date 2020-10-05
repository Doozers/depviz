// Package gomodhack ensures that `go mod` can detect some required dependencies.
// This package should not be imported directly.
package gomodhack

import (
	_ "github.com/gogo/protobuf/gogoproto"                                // required by protoc
	_ "github.com/gogo/protobuf/types"                                    // required by protoc
	_ "github.com/golang/protobuf/proto"                                  // nolint:staticcheck // we want this version; required by protoc
	_ "github.com/golang/protobuf/ptypes/timestamp"                       // required by protoc
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options" // required by protoc
)
