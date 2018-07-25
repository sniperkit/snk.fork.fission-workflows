#!/usr/bin/env bash


set -eu

. $(dirname $0)/utils.sh

NS=fission
NS_FUNCTION=fission-function
NS_BUILDER=fission-builder
fissionHelmId=fission
fissionWorkflowsHelmId=fission-workflows
TEST_STATUS=0
TEST_LOGFILE_PATH=tests.log
BIN_DIR="${BIN_DIR:-$HOME/testbin}"


cleanup_fission_workflows() {
    emph "Removing Fission Workflows deployment..."
    helm_uninstall_release ${fissionWorkflowsHelmId}
    # TODO cleanup workflow functions too
}

cleanup_fission() {
    cleanup_fission_workflows
    emph "Removing Fission deployment..."
    helm_uninstall_release ${fissionHelmId}

    emph "Removing custom resources..."
    clean_tpr_crd_resources || true
    kubectl delete all --all -n ${NS} --now
    kubectl delete all --all -n ${NS_FUNCTION} --now
    kubectl delete all --all -n ${NS_BUILDER} --now


    # Trigger deletion of all namespaces before waiting - for concurrency of deletion
    emph "Forcing deletion of namespaces..."
    kubectl delete ns/${NS} --now > /dev/null 2>&1 # Sometimes it is not deleted by helm delete
    kubectl delete ns/${NS_BUILDER} --now > /dev/null 2>&1 # Sometimes it is not deleted by helm delete
    kubectl delete ns/${NS_FUNCTION} --now > /dev/null 2>&1 # Sometimes it is not deleted by helm delete

    # Wait until all namespaces are actually deleted!
    emph "Awaiting deletion of namespaces..."
    verify_ns_deleted() {
        kubectl delete ns/${1} --now 2>&1  | grep -qv "Error from server (Conflict):"
    }
    # Namespaces sometimes take a long time to delete for some reason
    sleep 10
    RETRY_LIMIT=10 RETRY_DELAY=10 retry verify_ns_deleted ${NS_BUILDER}
    RETRY_LIMIT=10 RETRY_DELAY=10 retry verify_ns_deleted ${NS}
    RETRY_LIMIT=10 RETRY_DELAY=10 retry verify_ns_deleted ${NS_FUNCTION}
}

reset_fission_crd_resources() {
    NS_CRDS=${1:-default}
    echo "TODO reset fission"
    exit 1
    # TODO remove all functions, etc.
    reset_crd_resources
}

# Ensure printing of report
retry cleanup_fission