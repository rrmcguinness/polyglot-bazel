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

package example.service;

import examples.java.grpc.AuditGrpc;
import examples.java.pb.Model;
import io.grpc.stub.StreamObserver;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

public class AuditService extends AuditGrpc.AuditImplBase {
    private static final Log LOG = LogFactory.getLog(AuditService.class);
    @Override
    public StreamObserver<Model.AuditRecord> create(StreamObserver<Model.AuditResponse> responseObserver) {
        return new StreamObserver<Model.AuditRecord>() {
            @Override
            public void onNext(Model.AuditRecord value) {
                LOG.info(String.format("Received %s", value));
            }

            @Override
            public void onError(Throwable t) {
                LOG.error("Errpr processing: %s", t);
            }

            @Override
            public void onCompleted() {
                LOG.info("Complete");
            }
        };
    }
}
