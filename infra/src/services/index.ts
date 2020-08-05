import config from "../config";
import { createRendezvousService } from "./rendezvous";
import { createModelEnricher } from "./enricher";

const services = (config.rendezvous.services || []).map((service) => {
  return [createRendezvousService(service), createModelEnricher(service)];
});

export { services };
