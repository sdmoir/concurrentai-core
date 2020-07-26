import { createModelService } from "./model";

import config, { ServiceConfig } from "../config";

const models = config.rendezvous.services.reduce(
  (models: Record<string, any>, service: ServiceConfig) => {
    models[service.id] = service.models.map(createModelService);
    return models;
  },
  {}
);

export { models };
