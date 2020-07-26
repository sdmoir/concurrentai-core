import * as k8s from "@pulumi/kubernetes";

import config, { ServiceConfig } from "../config";
import { provider } from "../cluster/provider";
import { secret as registrySecret } from "../cluster/registry";

export function createRendezvousService(serviceConfig: ServiceConfig) {
  const metadata = { name: `rendezvous-service-${serviceConfig.id}` };
  const appLabels = { run: `rendezvous-service-${serviceConfig.id}` };

  const deployment = new k8s.apps.v1.Deployment(
    `rendezvous-deployment-${serviceConfig.id}`,
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
                name: "rendezvous-api",
                ports: [{ containerPort: 9000 }],
                image: `registry.digitalocean.com/concurrent-ai/rendezvous-service-poc-api:latest`,
                imagePullPolicy: "Always",
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
                image: `registry.digitalocean.com/concurrent-ai/rendezvous-service-poc-collector:latest`,
                imagePullPolicy: "Always",
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
                  {
                    name: "MODEL_ENDPOINT",
                    value: `http://rendezvous-model-${
                      serviceConfig.models[0]?.id || "active"
                    }/invocations`,
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
            imagePullSecrets: [
              {
                name: registrySecret.metadata.name,
              },
            ],
          },
        },
      },
    },
    {
      provider,
    }
  );

  const service = new k8s.core.v1.Service(
    `rendezvous-service-${serviceConfig.id}`,
    {
      metadata: metadata,
      spec: {
        ports: [{ port: 80, targetPort: 9000 }],
        selector: appLabels,
      },
    },
    {
      provider,
    }
  );

  return {
    deployment,
    service,
  };
}
