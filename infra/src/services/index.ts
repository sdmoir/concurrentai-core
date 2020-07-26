import { createRendezvousService } from "./rendezvous";

import config from "../config";

const services = (config.rendezvous.services || []).map(
  createRendezvousService
);

export { services };
