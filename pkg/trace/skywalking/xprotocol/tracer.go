/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package xprotocol

import (
	"context"
	"mosn.io/mosn/pkg/protocol/xprotocol"
	"strconv"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	"github.com/SkyAPM/go2sky/reporter/grpc/common"
	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/protocol/http"
	"mosn.io/mosn/pkg/trace"
	"mosn.io/mosn/pkg/trace/skywalking"
	"mosn.io/mosn/pkg/types"
)


func init() {
	trace.RegisterTracerBuilder(skywalking.SkyDriverName, protocol.Xprotocol, NewHttpSkyTracer)
}

type xprotoSkyTracer struct {
	*go2sky.Tracer
}

func NewHttpSkyTracer(_ map[string]interface{}) (types.Tracer, error) {
	return &xprotoSkyTracer{}, nil
}

func (tracer *xprotoSkyTracer) SetGO2SkyTracer(t *go2sky.Tracer) {
	tracer.Tracer = t
}

func (tracer *xprotoSkyTracer) Start(ctx context.Context, frame interface{}, _ time.Time) types.Span {
	xframe, ok := frame.(xprotocol.XFrame)
	if !ok || xframe == nil {
		return skywalking.NoopSpan
	}

	// ignore heartbeat
	if xframe.IsHeartbeatFrame() {
		return skywalking.NoopSpan
	}

	// use delegate instrument if exists
	// subProtocol := types.ProtocolName(mosnctx.Get(ctx, types.ContextSubProtocol).(string))

	return httpSkySpan{
		tracer: tracer,
		ctx:    context.TODO(), // TODO
		carrier: &spanCarrier{
			entrySpan: nil, // TODO
		},
	}
}

type spanCarrier struct {
	entrySpan go2sky.Span
	exitSpan  go2sky.Span
}

type httpSkySpan struct {
	skywalking.SkySpan
	tracer  *xprotoSkyTracer
	ctx     context.Context
	carrier *spanCarrier
}

func (h httpSkySpan) TraceId() string {
	return go2sky.TraceID(h.ctx)
}

func (h httpSkySpan) InjectContext(requestHeaders types.HeaderMap, requestInfo api.RequestInfo) {
	header, ok := requestHeaders.(http.RequestHeader)
	if !ok || header.RequestHeader == nil {
		log.DefaultLogger.Debugf("[SkyWalking] [tracer] [http1] unable to get request header, upstream trace ignored")
		return
	}
	requestURI := string(header.RequestURI())
	upstreamLocalAddress := requestInfo.UpstreamLocalAddress()

	// create exit span (upstream)
	exit, err := h.tracer.CreateExitSpan(h.ctx, requestURI, upstreamLocalAddress, func(header string) error {
		requestHeaders.Add(propagation.Header, header)
		return nil
	})
	if err != nil {
		log.DefaultLogger.Errorf("[SkyWalking] [tracer] [http1] create exit span error, err: %v", err)
		return
	}

	exit.SetComponent(skywalking.MOSNComponentID)
	exit.SetSpanLayer(common.SpanLayer_Http)
	h.carrier.exitSpan = exit
}

func (h httpSkySpan) SetRequestInfo(requestInfo api.RequestInfo) {
	responseCode := strconv.Itoa(requestInfo.ResponseCode())

	// end exit span (upstream)
	if h.carrier.exitSpan != nil {
		exit := h.carrier.exitSpan
		if requestInfo.ResponseCode() >= http.BadRequest {
			exit.Error(time.Now(), skywalking.ErrorLog)
		}
		exit.Tag(go2sky.TagStatusCode, responseCode)
		exit.End()
	}

	// entry span (downstream)
	entry := h.carrier.entrySpan
	if requestInfo.ResponseCode() >= http.BadRequest {
		entry.Error(time.Now(), skywalking.ErrorLog)
	}
	entry.Tag(go2sky.TagStatusCode, responseCode)
	// TODO More mosn information
}

func (h httpSkySpan) FinishSpan() {
	h.carrier.entrySpan.End()
}
