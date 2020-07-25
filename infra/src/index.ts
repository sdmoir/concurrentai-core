import config, { ServiceConfig } from "./config";
import { createRendezvousGateway } from "./deployments/krakend";
import { createRendezvousService } from "./deployments/rendezvous";
import { createModelService } from "./deployments/model";
import { DeployableService } from "./deployments/types";

export const gateway = createRendezvousGateway(config);

export const services = config.rendezvous.services.map((service) =>
  createRendezvousService(config, service)
);

export const models = config.rendezvous.services.reduce(
  (models: Record<string, DeployableService[]>, service: ServiceConfig) => {
    models[service.id] = service.models.map((model) =>
      createModelService(config, model)
    );
    return models;
  },
  {}
);
