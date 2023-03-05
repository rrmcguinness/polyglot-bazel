/*
 * Copyright 2023 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package example.client;

import examples.java.grpc.EventsGrpc;
import examples.java.pb.Model;
import io.grpc.Channel;
import io.grpc.stub.StreamObserver;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

public class EventClient {
    private static final Log LOG = LogFactory.getLog(EventClient.class);

    private final EventsGrpc.EventsBlockingStub blockingStub;
    private final EventsGrpc.EventsStub asyncStub;

    public EventClient(Channel channel) {
        blockingStub = EventsGrpc.newBlockingStub(channel);
        asyncStub = EventsGrpc.newStub(channel);
    }

    public List<Model.EventResponse> put(Model.Event... events) throws InterruptedException {
        final List<Model.EventResponse> responses = new ArrayList<>();
        final CountDownLatch latch = new CountDownLatch(1);
        StreamObserver<Model.EventResponse> responseObserver = new StreamObserver<Model.EventResponse>() {
            @Override
            public void onNext(Model.EventResponse value) {
                responses.add(value);
            }

            @Override
            public void onError(Throwable t) {
                LOG.error(t);
                latch.countDown();
            }

            @Override
            public void onCompleted() {
                LOG.info("finished put response listener");
                latch.countDown();
            }
        };

        StreamObserver<Model.Event> eventStreamObserver = asyncStub.put(responseObserver);

        for (Model.Event event : events) {
            eventStreamObserver.onNext(event);
        }
        eventStreamObserver.onCompleted();

        if (!latch.await(1, TimeUnit.SECONDS)) {
            LOG.warn("route incomplete waiting 1 second");
        }

        return responses;
    }

}
