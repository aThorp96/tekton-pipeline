ARG GO_BUILDER=brew.registry.redhat.io/rh-osbs/openshift-golang-builder:v1.23
ARG RUNTIME=registry.access.redhat.com/ubi9/ubi-minimal:latest@sha256:92b1d5747a93608b6adb64dfd54515c3c5a360802db4706765ff3d8470df6290

FROM $RUNTIME as dependency-builder

COPY dependencies/tini dependencies/tini
WORKDIR /dependencies/tini
RUN microdnf update && microdnf install -y cmake gcc
ENV CFLAGS="-DPR_SET_CHILD_SUBREAPER=36 -DPR_GET_CHILD_SUBREAPER=37"
RUN cmake . && make tini

FROM $GO_BUILDER AS builder

WORKDIR /go/src/github.com/tektoncd/pipeline
COPY upstream .
COPY .konflux/patches patches/
RUN set -e; for f in patches/*.patch; do echo ${f}; [[ -f ${f} ]] || continue; git apply ${f}; done
COPY head HEAD
ENV GODEBUG="http2server=0"
RUN go build -ldflags="-X 'knative.dev/pkg/changeset.rev=$(cat HEAD)'" -mod=vendor -tags disable_gcp -v -o /tmp/resolvers \
    ./cmd/resolvers

FROM $RUNTIME
ARG VERSION=pipeline-main

ENV RESOLVERS=/usr/local/bin/resolvers \
    KO_APP=/ko-app \
    KO_DATA_PATH=/kodata

COPY --from=builder /tmp/resolvers /ko-app/resolvers
COPY head ${KO_DATA_PATH}/HEAD

COPY --from=dependency-builder /dependencies/tini/tini /sbin/tini
RUN chmod 0755 /sbin/tini && chown root:root /sbin/tini

LABEL \
      com.redhat.component="openshift-pipelines-resolvers-rhel8-container" \
      name="openshift-pipelines/pipelines-resolvers-rhel8" \
      version=$VERSION \
      summary="Red Hat OpenShift Pipelines Resolvers" \
      maintainer="pipelines-extcomm@redhat.com" \
      description="Red Hat OpenShift Pipelines Resolvers" \
      io.k8s.display-name="Red Hat OpenShift Pipelines Resolvers" \
      io.k8s.description="Red Hat OpenShift Pipelines Resolvers" \
      io.openshift.tags="pipelines,tekton,openshift"

RUN microdnf install -y shadow-utils
RUN groupadd -r -g 65532 nonroot && useradd --no-log-init -r -u 65532 -g nonroot nonroot
USER 65532

ENTRYPOINT ["/sbin/tini", "--", "/ko-app/resolvers"]
