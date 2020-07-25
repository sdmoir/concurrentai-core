import * as k8s from "@pulumi/kubernetes";

export interface DeployableService {
  deployment: k8s.apps.v1.Deployment;
  service: k8s.core.v1.Service;
}
