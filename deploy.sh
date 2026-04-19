#!/bin/bash
eval $(minikube docker-env)
docker build -t lol ./src
kubectl apply -f kubik/configmap.yaml

kubectl apply -f kubik/pod.yaml
kubectl wait --for=condition=ready pod/lolpod --timeout=60s

kubectl apply -f kubik/deployment.yaml
kubectl rollout restart deployment lol

kubectl apply -f kubik/service.yaml
kubectl apply -f kubik/daemonset.yaml
kubectl apply -f kubik/cronjob.yaml
kubectl apply -f kubik/statefulset.yaml

kubectl rollout status statefulset/lol-stateful
kubectl rollout status deployment/lol

istioctl install --set profile=demo -y
kubectl label namespace default istio-injection=enabled --overwrite
kubectl rollout restart deployment lol
kubectl rollout restart statefulset lol-stateful
 
kubectl apply -f mash/gateway.yaml
kubectl apply -f mash/virtualservice.yaml
kubectl apply -f mash/destinationrule.yaml
kubectl rollout status deployment/istio-ingressgateway -n istio-system --timeout=120s