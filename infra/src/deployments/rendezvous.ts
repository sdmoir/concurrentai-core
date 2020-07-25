import * as k8s from "@pulumi/kubernetes";

import { ServiceConfig, InfraConfig } from "../config";
import { DeployableService } from "./types";

export function createRendezvousService(
  config: InfraConfig,
  serviceConfig: ServiceConfig
): DeployableService {
  const metadata = { name: "rendezvous-service" };
  const appLabels = { run: "rendezvous-service" };

  const deployment = new k8s.apps.v1.Deployment("rendezvous-deployment", {
    metadata: metadata,
    spec: {
      selector: { matchLabels: appLabels },
      replicas: 1,
      template: {
        metadata: { labels: appLabels },
        spec: {
          containers: [
            {
              name: "rendezvous-api",
              image: `registry.digitalocean.com/concurrent-ai/rendezvous-${serviceConfig.id}-api:${process.env.GITHUB_SHA}`,
              env: [
                {
                  name: "KAFKA_BROKERS",
                  value: config.kafka.brokers,
                },
                {
                  name: "KAFKA_API_KEY",
                  value: config.kafka.apiKey,
                },
                {
                  name: "KAFKA_API_SECRET",
                  value: config.kafka.apiSecret,
                },
                {
                  name: "KAFKA_TOPIC",
                  value: serviceConfig.businessTopic,
                },
              ],
              volumeMounts: [
                {
                  name: "rendezvous-sockets",
                  mountPath: "/sockets",
                },
              ],
            },
            {
              name: "rendezvous-collector",
              image: `registry.digitalocean.com/concurrent-ai/rendezvous-${serviceConfig.id}-collector:${process.env.GITHUB_SHA}`,
              env: [
                {
                  name: "KAFKA_BROKERS",
                  value: config.kafka.brokers,
                },
                {
                  name: "KAFKA_API_KEY",
                  value: config.kafka.apiKey,
                },
                {
                  name: "KAFKA_API_SECRET",
                  value: config.kafka.apiSecret,
                },
                {
                  name: "KAFKA_TOPIC",
                  value: serviceConfig.collectionTopic,
                },
              ],
              volumeMounts: [
                {
                  name: "rendezvous-sockets",
                  mountPath: "/sockets",
                },
              ],
            },
          ],
          volumes: [
            {
              name: "rendezvous-sockets",
              emptyDir: {},
            },
          ],
        },
      },
    },
  });

  const service = new k8s.core.v1.Service("rendezvous-service", {
    metadata: metadata,
    spec: {
      ports: [{ port: 80, targetPort: 9000 }],
      selector: appLabels,
    },
  });

  return {
    deployment,
    service,
  };
}
