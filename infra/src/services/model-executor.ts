import * as k8s from "@pulumi/kubernetes";

import config, { ServiceConfig } from "../config";
import { provider } from "../cluster/provider";
import { secret as registrySecret } from "../cluster/registry";

export function createModelExecutor(serviceConfig: ServiceConfig) {
  const metadata = { name: `rendezvous-${serviceConfig.id}-model-executor` };
  const appLabels = { run: `rendezvous-${serviceConfig.id}-model-executor` };

  const deployment = new k8s.apps.v1.Deployment(
    `rendezvous-${serviceConfig.id}-model-executor-deployment`,
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
                name: "model-executor",
                image: `registry.digitalocean.com/concurrent-ai/rendezvous-model-executor:latest`,
                imagePullPolicy: "Always",
                env: [
                  {
                    name: "ORGANIZATION_ID",
                    value: config.rendezvous.organizationId,
                  },
                  {
                    name: "SERVICE_ID",
                    value: serviceConfig.id,
                  },
                  {
                    name: "PULSAR_URL",
                    value: config.pulsar.url,
                  },
                  {
                    name: "MODEL_ENDPOINT",
                    value: `http://rendezvous-${serviceConfig.id}-model-${
                      serviceConfig.models[0]?.id || "active"
                    }/invocations`,
                  },
                ],
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

  return {
    deployment,
  };
}
