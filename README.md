# Prometheus-Grafana

This is a **Work-In-Progress** Repo for learning how monitor your kubernetes clusters using prometheus and visualize using grafana.

Both Prometheus and Grafana have a very good documentation 

[Prometheus Docs]("https://prometheus.io/docs/introduction/overview/")

[Grafana Docs]("https://grafana.com/docs/grafana/latest/")

## Pre-Requisite

- Kubernetes Cluster (can be minikube)
- Helm 

If you don't have them installed. Follow the below links:

[Install Minikube]("https://minikube.sigs.k8s.io/docs/start/")

[Install Helm]("https://helm.sh/docs/intro/install/")

<br></br>

## It all starts with Monitoring 

Monitoring your Kubernetes cluster is essential for ensuring the health and performance of your applications and infrastructure. Here are some reasons why monitoring your Kubernetes cluster is important:

- Identify issues and troubleshoot: By monitoring your Kubernetes cluster, you can quickly identify issues such as application crashes, resource bottlenecks, and network problems. With real-time monitoring, you can troubleshoot issues before they escalate and impact your users.

- Optimize performance and capacity: Monitoring allows you to track the performance of your applications and infrastructure over time, and identify opportunities to optimize performance and capacity. By understanding usage patterns and resource consumption, you can make informed decisions about scaling your infrastructure and improving the efficiency of your applications.

- Ensure high availability: Kubernetes is designed to provide high availability for your applications, but this requires careful monitoring and management. By monitoring your cluster and setting up alerts, you can ensure that your applications remain available even in the event of failures or unexpected events.

- Security and compliance: Monitoring your Kubernetes cluster can help you identify potential security risks and ensure compliance with regulations and policies. By tracking access logs and other security-related metrics, you can quickly detect and respond to potential security threats.

<br></br>

## Using Prometheus for monitoring

Prometheus is an open-source monitoring and alerting system that helps you collect and store metrics about your software systems and infrastructure, and analyze that data to gain insights into their health and performance. It provides a powerful query language, a flexible data model, and a range of integrations with other tools and systems. With Prometheus, you can easily monitor metrics such as CPU usage, memory usage, network traffic, and application-specific metrics, and use that data to troubleshoot issues, optimize performance, and create alerts to notify you when things go wrong.

<br></br>

## Why Prometheus over other monitoring tools ?

Prometheus is a popular choice for Kubernetes monitoring for several reasons:

- Open-source: Prometheus is an open-source project that is free to use and has a large community of contributors. This means that you can benefit from ongoing development, bug fixes, and feature enhancements without paying for a commercial monitoring solution.

- Native Kubernetes support: Prometheus is designed to work seamlessly with Kubernetes, making it easy to deploy and integrate with your Kubernetes environment. It provides pre-configured Kubernetes dashboards and supports auto-discovery of Kubernetes services and pods.

- Powerful query language: Prometheus provides a powerful query language that allows you to easily retrieve and analyze metrics data. This allows you to create custom dashboards and alerts, and to troubleshoot issues more easily.

- Scalability: Prometheus is designed to be highly scalable, allowing you to monitor large and complex Kubernetes environments with ease. It supports multi-node architectures and can handle large volumes of data without significant performance degradation.

- Integrations: Prometheus integrates with a wide range of other tools and systems, including Grafana for visualization, Alertmanager for alerting, and Kubernetes API server for metadata discovery.

<br></br>

## Prometheus Architecture

<br></br>
![Alt text](https://prometheus.io/assets/architecture.png)

<br></br>

## What is Grafana ?

Grafana is a popular open-source data visualization and analytics platform that allows you to create custom dashboards and visualizations based on a variety of data sources. Grafana is often used for monitoring and analyzing metrics and logs in real-time, making it an ideal tool for monitoring systems and applications, including Kubernetes environments.

Grafana supports a wide range of data sources, including databases, time-series databases, and other data storage systems. It provides a powerful query language that allows you to retrieve and analyze data from these sources, and to create custom dashboards and alerts based on that data.

In addition to its powerful data visualization and analysis capabilities, Grafana is also highly extensible. It supports a wide range of plugins and integrations, including integrations with popular monitoring and logging tools like Prometheus, Elasticsearch, and InfluxDB.

# Lets dive into Practical part 

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/99ccf4db-1b17-4ae9-8b6c-6434e8989b21)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/f58dfc08-e43d-4937-ba69-307fb6c7888d)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/c5c45315-89b4-4e9a-844a-03bb90e8e568)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/cb854f60-9f88-460d-bd0f-935a3246140b)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/fc34a791-85a6-4418-bfde-35cc10d2a594)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/ee2a4f2e-4fe6-4331-a1d7-2dc325a4694b)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/ea237810-948d-43bb-b329-0e5ac12836c2)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/f3a901a9-d004-4dd1-bdcb-7165540f95b4)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/b72de8c4-6bd3-42e6-81ba-658977f5f91a)

# Prometheus running

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/c3112d80-a236-4e0b-a160-f120e60ffcaf)

# go to  Grafana  - 

# go to  Grafana  - https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/blob/main/Installation/Grafana/helm.md

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/e2b7de19-61af-47ac-b37f-02e1658b7262)

# Copy the Cmd step by step

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/db5a40e8-7178-4fef-b758-7c21dd035850)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/7c5d078f-3c7b-4e23-8bb2-774565251d41)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/bd58d024-ffa9-48fe-ab53-4294bb6bb8c1)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/2306b088-16fb-4c9b-bbdf-205603189e21)

# Login to Grafana + password …

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/5486793c-59b3-49cf-ac00-f97884c5089c)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/6a7e5eb2-b882-4108-a706-ee3a3f08d260)

# Expose Grafana with New Service

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/893901bd-2f04-4dba-aeea-961964dbfba9)

# Kubectl get svc

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/5fbacc89-5338-4b74-95a8-0fe80e021a5c)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/b62404c1-fcef-45c6-b9d8-56d4d72fa261)

# 31281/TCP  Check in the host server

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/ea440e20-b717-4f67-8e3c-900bacd89b1e)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/0b8aef7a-18bf-440d-bd5f-1c5521333325)

# Create Prometheus as data source for Grafana

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/f8d143de-5c27-4eee-bcbc-7766fb7d4f43)

# Add Data source

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/77fc618f-a56b-42e7-8ae0-fe6afea26cc8)

# Provide IP address of Prometheus  + save it  + test it 

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/25f0ebe8-3f35-4cc4-add8-b8aa84936cac)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/687bb565-071f-4396-a488-8a6a7c6c995d)

# Data Source is Working

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/b50d8b6f-6128-4c96-a317-ddecfe4b6ffc)

# now the Data from Prometheus will be retrived in Grafana 

# now See the dashboard

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/5f6cf99d-3417-4281-97f6-e864503f375c)

# Click on dashboard + import

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/29fa8ee4-13ff-4cd8-81b8-2b252bb8f946)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/d6a4a109-df7d-4adc-957d-75c34848203f)

# ID 3662 + click on load

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/4020967f-75f9-4002-b960-4710ba0542e4)

# Click on import

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/4be6db80-2e7c-41b6-bbed-7e7790b07417)

# Now the below dashboard extracting data from Minikube cluster

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/00d63ab7-712f-440c-9015-0cc2a9a7c792)

# Now make new Entry in the list & see the change

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/a62634fd-ff03-490b-b6eb-a590707e11af)

# Run Kubectl get svc

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/80ba35dd-a582-44d5-b71c-936579a4da58)

# Kube-state-metric Running on port 30421./TCP .. Now  see the magic of Grafana

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/fa9c199a-e6c8-4785-96a3-9942821e102a)

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/95b395ff-7b47-45fe-b727-f983210de47c)

# click on metrics

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/4f09f108-b294-4f60-a0af-ebe7ac19423c)

# See the status of Deployment

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/ec8ab8e0-826a-47b8-92ea-22f1ff0584b8)

# Status check in Prometheus

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/7ce96be7-f6a9-4cc8-9685-2b0188d1ecef)

# now See the beauty of Grafana it will  display the above data source(Prometheus) in a visualization format

# Reload in Grafana

![image](https://github.com/Siddhartha082/Kubernetes_Monitoring_Prometheus-Grafana-Tutorial/assets/110781138/4560009c-6bfb-4978-aa0e-62c58bb2537d)











































