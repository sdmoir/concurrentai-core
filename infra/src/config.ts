import * as pulumi from "@pulumi/pulumi";

interface KafkaConfig {
  clientId: string;
  brokers: string;
  apiKey: string;
  apiSecret: string;
}

export interface ModelConfig {
  id: string;
  image: string;
}

export interface ServiceConfig {
  id: string;
  businessTopic: string;
  collectionTopic: string;
  models: [ModelConfig];
}

export interface RendezvousConfig {
  organizationId: string;
  services: [ServiceConfig];
}

export interface InfraConfig {
  kafka: KafkaConfig;
  rendezvous: RendezvousConfig;
}

const config = new pulumi.Config();

const infraConfig: InfraConfig = {
  kafka: config.requireObject<KafkaConfig>("kafka"),
  rendezvous: config.requireObject<RendezvousConfig>("rendezvous"),
};

export default infraConfig;
