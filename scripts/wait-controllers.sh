#!/usr/bin/env bash

if [[ ${1} -eq "dependencies" ]]; then
  kubectl rollout status deployment -n cert-manager cert-manager
  kubectl rollout status deployment -n cert-manager cert-manager-cainjector
  kubectl rollout status deployment -n cert-manager cert-manager-webhook

elif [[ ${1} -eq "cluster-api" ]]; then
  kubectl rollout status deployment -n capd-system capd-controller-manager
  kubectl rollout status deployment -n capi-kubeadm-bootstrap-system capi-kubeadm-bootstrap-controller-manager
  kubectl rollout status deployment -n capi-kubeadm-control-plane-system capi-kubeadm-control-plane-controller-manager
  kubectl rollout status deployment -n capi-system capi-controller-manager
fi

sleep 60