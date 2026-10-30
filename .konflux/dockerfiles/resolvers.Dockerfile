ARG GO_BUILDER=registry.access.redhat.com/hi/go:1.27-fips-builder@sha256:b7a89799fe7ac875c0923329ea6cdcb61837027b5a44fa9f23a88101c05e833a
ARG UBI_RUNTIME=registry.redhat.io/ubi9/ubi-minimal-pqc:latest@sha256:72187ec58a94bd958ff48e41881d7c5d8bc0b49586cc1661da5ba4d92374c206
ARG GIT_BASE=registry.access.redhat.com/hi/git:latest

FROM $UBI_RUNTIME AS dependency-builder

COPY dependencies/tini dependencies/tini
WORKDIR /dependencies/tini
RUN microdnf update -y && microdnf install -y cmake gcc
ENV CFLAGS="-DPR_SET_CHILD_SUBREAPER=36 -DPR_GET_CHILD_SUBREAPER=37"
RUN cmake . && make tini
RUN chmod 0755 /dependencies/tini/tini && chown root:root /dependencies/tini/tini

FROM $GO_BUILDER AS builder

WORKDIR /go/src/github.com/tektoncd/pipeline
COPY upstream .
COPY .konflux/patches patches/
RUN set -e; for f in patches/*.patch; do echo ${f}; [[ -f ${f} ]] || continue; git apply ${f}; done
COPY head HEAD
ENV GODEBUG="http2server=0"
RUN go build -ldflags="-X 'knative.dev/pkg/changeset.rev=$(cat HEAD)'" -mod=vendor -tags disable_gcp,strictfipsruntime -v -o /tmp/resolvers \
    ./cmd/resolvers

FROM $GIT_BASE
ARG VERSION=next

ENV RESOLVERS=/usr/local/bin/resolvers \
    KO_APP=/ko-app \
    KO_DATA_PATH=/kodata

COPY --from=builder /tmp/resolvers /ko-app/resolvers
COPY head ${KO_DATA_PATH}/HEAD

COPY --from=dependency-builder /dependencies/tini/tini /sbin/tini

LABEL \
    com.redhat.component="openshift-pipelines-resolvers-rhel9-container" \
    cpe="cpe:/a:redhat:openshift_pipelines:next::el9" \
    description="Red Hat OpenShift Pipelines tektoncd-pipeline resolvers" \
    io.k8s.description="Red Hat OpenShift Pipelines tektoncd-pipeline resolvers" \
    io.k8s.display-name="Red Hat OpenShift Pipelines tektoncd-pipeline resolvers" \
    io.openshift.tags="tekton,openshift,tektoncd-pipeline,resolvers" \
    maintainer="pipelines-extcomm@redhat.com" \
    name="openshift-pipelines/pipelines-resolvers-rhel9" \
    summary="Red Hat OpenShift Pipelines tektoncd-pipeline resolvers" \
    version="next"

USER 65532:65532

ENTRYPOINT ["/sbin/tini", "--", "/ko-app/resolvers"]

