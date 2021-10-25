# Introduction

The Gilman Lab (GLab) project is aimed at creating a self-service cloud-like
architecture that can be run in a homelab environment for learning purposes. 
While there are dozens of ways to create such an architecture in today's world,
GLab is opionated towards using a non-Kubernetes based approach which focuses
on leveraging a number of Hashicorp and other third-party open source software
in order to create the final product.

One of the primary goals of GLab is that of reproduceability. Since the
project is aimed creating a learning environment where modern operational
techniques can be tested and improved - it's very likely for the environment to
end up in an incompatible state on a frequent basis. As such, being able to
quickly recover and rebuild the entire architecture is key. 

The GLab project aims to use open-source software where applicable. The only
caveat to this is the hypervisor which relies on a vCenter cluster with multiple 
ESXi nodes. The project aims to abstract away the hypervisor as much as possible
by relying heavily on packing most of the functionality into the individual
virtual machines - thus enabling it to be run anyhwere VMs can run.

The latter is important to the last goal of GLab which is portability. To
increase productivity and make iterating faster - the entire GLab stack can and
should be run locally on a development machine. A fresh instance can be spun up
locally for development and testing purposes.

# The Stack

The GLab project is made up of a core stack of software which makes up the
self-service cloud architecture:

* Flatcar Linux
  * Provides the operating systems that all nodes run on
  * The Ignition file provides a common way for configuring nodes
* Docker
  * Provides the container runtime in which services operate in
* Portworx
  * Serves as the primary storage provider
  * Provides distributed volumes for use by Docker
* Consul
  * Serves as the backend for Vault and Nomad
  * Provides configuration data to Nomad and other services via the KV store
  * Provides a secure service mesh for inter-service communication
* Vault
  * Serves as the all-in-one secret management service
  * Provides storage for plain-text secrets via the KV store
  * Acts as the central certificate authority for mTLS and SSH
  * Provides transit-based encryption services for various processes
* Nomad
  * Serves as the primary orchestration service
  * Provides an interface for configuring and bringing up services
  * Acts as the glue for Consul and Vault
* Keycloak
  * Serves as the primary identity platform for authentication and authorization
  * Integrates with other services for SSO-like support
  * Provides an interface for permission based control

In addition to the above services, the following tools are also utilized in
bringing up the stack:

* Vagrant
  * Provides the abstraction layer for running the stack locally for development
* Ansible
  * Provides the primary mechanism for applying configuration to nodes
  * Acts as an authoritative source for node and service configuration