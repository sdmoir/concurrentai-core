import * as k8s from "@pulumi/kubernetes";

import { ModelConfig, ServiceConfig } from "../config";
import { provider } from "../cluster/provider";
import { secret as registrySecret } from "../cluster/registry";

export function createModelService(
  serviceConfig: ServiceConfig,
  modelConfig: ModelConfig
) {
  const fullModelId = `${serviceConfig.id}-${modelConfig.id}`;
  const metadata = { name: `rendezvous-model-${fullModelId}` };
  const appLabels = { run: `rendezvous-model-${fullModelId}` };

  const deployment = new k8s.apps.v1.Deployment(
    `rendezvous-model-${fullModelId}-deployment`,
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
                name: "model",
                ports: [{ containerPort: 8080 }],
                image: modelConfig.image,
                imagePullPolicy: "Always",
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
    `rendezvous-model-${fullModelId}-service`,
    {
      metadata: metadata,
      spec: {
        ports: [{ port: 80, targetPort: 8080 }],
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
