import * as k8s from "@pulumi/kubernetes";
import { InfraConfig } from "../config";
import { DeployableService } from "./types";

export function createRendezvousGateway(
  config: InfraConfig
): DeployableService {
  const metadata = { name: "rendezvous-gateway" };
  const appLabels = { app: "krakend" };

  const deployment = new k8s.apps.v1.Deployment(
    "rendezvous-gateway-deployment",
    {
      metadata: metadata,
      spec: {
        selector: { matchLabels: appLabels },
        replicas: 1,
        template: {
          metadata: { labels: appLabels },
          spec: {
            containers: [
              {
                name: "krakend",
                image: `registry.digitalocean.com/concurrent-ai/rendezvous-krakend:${process.env.GITHUB_SHA}`,
                imagePullPolicy: "Always",
                ports: [{ containerPort: 8080 }],
                command: ["/usr/bin/krakend"],
                args: [
                  "run",
                  "-d",
                  "-c",
                  "/etc/krakend/krakend.json",
                  "-p",
                  "8080",
                ],
              },
            ],
          },
        },
      },
    }
  );

  const service = new k8s.core.v1.Service("rendezvous-gateway-service", {
    metadata: metadata,
    spec: {
      type: "LoadBalancer",
      ports: [{ name: "http", port: 80, targetPort: 8080, protocol: "TCP" }],
      selector: appLabels,
    },
  });

  return {
    deployment,
    service,
  };
}
