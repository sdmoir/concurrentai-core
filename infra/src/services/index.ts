import config from "../config";
import { createRendezvousService } from "./rendezvous";
import { createModelEnricher } from "./enricher";
import { createModelExecutor } from "./model-executor";

const services = (config.rendezvous.services || []).map((service) => {
  return [
    createRendezvousService(service),
    createModelEnricher(service),
    createModelExecutor(service),
  ];
});

export { services };
