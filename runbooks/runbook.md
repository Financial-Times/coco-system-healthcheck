# UPP - System Healthcheck

A service that monitors key system metrics for the EC2 instances in our containerised clusters.

## Code

upp-system-healthcheck

## Primary URL

https://github.com/Financial-Times/coco-system-healthcheck

## Service Tier

Bronze

## Lifecycle Stage

Production

## Delivered By

content

## Supported By

content

## Known About By

- dimitar.terziev
- elitsa.pavlova
- hristo.georgiev
- donislav.belev
- mihail.mihaylov
- boyko.boykov

## Host Platform

AWS

## Architecture

A service that monitors key system metrics for the EC2 instances in our containerised clusters. See the service GitHub repository for more details.

## Contains Personal Data

No

## Contains Sensitive Data

No

## Failover Architecture Type

ActiveActive

## Failover Process Type

FullyAutomated

## Failback Process Type

FullyAutomated

## Failover Details

The service is deployed in all clusters as a daemonset.

## Data Recovery Process Type

NotApplicable

## Data Recovery Details

The service does not store data, so it does not require any data recovery steps.

## Release Process Type

PartiallyAutomated

## Rollback Process Type

Manual

## Release Details

The release is triggered by making a Github release which is then picked up by a Jenkins multibranch pipeline. The Jenkins pipeline should be manually started in order for it to deploy the helm package to the Kubernetes clusters.

## Key Management Process Type

NotApplicable

## Key Management Details

There is no key rotation procedure for this system.

## Monitoring

Pod health:

- <https://upp-prod-delivery-eu.ft.com/__health/__pods-health?service-name=system-healthcheck>
- <https://upp-prod-publish-eu.ft.com/__health/__pods-health?service-name=system-healthcheck>
- <https://upp-prod-delivery-us.ft.com/__health/__pods-health?service-name=system-healthcheck>
- <https://upp-prod-publish-us.ft.com/__health/__pods-health?service-name=system-healthcheck>
- <https://pac-prod-eu.upp.ft.com/__health/__pods-health?service-name=system-healthcheck>
- <https://pac-prod-us.upp.ft.com/__health/__pods-health?service-name=system-healthcheck>

## First Line Troubleshooting

[First Line Troubleshooting guide](https://github.com/Financial-Times/upp-docs/tree/master/guides/ops/first-line-troubleshooting)

## Second Line Troubleshooting

Please refer to the GitHub repository README for troubleshooting information.
