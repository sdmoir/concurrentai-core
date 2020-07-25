import * as k8s from "@pulumi/kubernetes";
import { ModelConfig, InfraConfig } from "../config";
import { DeployableService } from "./types";

export function createModelService(
  config: InfraConfig,
  modelConfig: ModelConfig
): DeployableService {
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
                image: modelConfig.image,
              },
            ],
          },
        },
      },
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
    }
  );

  return {
    deployment,
    service,
  };
}
