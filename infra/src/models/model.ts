import * as k8s from "@pulumi/kubernetes";

import { ModelConfig } from "../config";
import { provider } from "../cluster/provider";
import { secret as registrySecret } from "../cluster/registry";

export function createModelService(modelConfig: ModelConfig) {
  const metadata = { name: `rendezvous-model-${modelConfig.id}` };
  const appLabels = { run: `rendezvous-model-${modelConfig.id}` };

  const deployment = new k8s.apps.v1.Deployment(
    `rendezvous-model-${modelConfig.id}-deployment`,
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
    `rendezvous-model-${modelConfig.id}-service`,
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
