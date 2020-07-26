import * as pulumi from "@pulumi/pulumi";

interface DigitalOceanConfig {
  registryToken: string;
}

interface KafkaConfig {
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
  digitalocean: DigitalOceanConfig;
  kafka: KafkaConfig;
  rendezvous: RendezvousConfig;
}

const config = new pulumi.Config();

const infraConfig: InfraConfig = {
  digitalocean: config.requireObject<DigitalOceanConfig>("digitalocean"),
  kafka: config.requireObject<KafkaConfig>("kafka"),
  rendezvous: config.requireObject<RendezvousConfig>("rendezvous"),
};

export default infraConfig;
