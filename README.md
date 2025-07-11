# Ogugu
Ogugu is a backend application for an RSS feed aggregator and reader. It allows users to add RSS feed URLs to a centralized database, subscribe to multiple feeds, and regularly fetch and aggregate new posts from these feeds for seamless consumption.

## Features
- Adding and managing RSS feed links in a shared database
- Subscribing to various RSS feeds to personalize content
- Periodic fetching and aggregation of RSS feed posts to keep user content up-to-date

### Built with
- Golang
- Postgres
- Redis
- Docker

## Prerequisite
1. [Golang 1.24+](https://go.dev/doc/install)
2. Docker — For containerized services (PostgreSQL & Redis)
3. Git
4. [Golang-Migrate](https://github.com/golang-migrate/migrate?tab=readme-ov-file)
5. [Swag Cli](https://github.com/swaggo/swag) - For API testing
### Infrastructure (Optional)
- kind
- kubectl
- terraform

## Setup
1. Clone the repository and navigate into the project directory
```bash
git clone https://github.com/tonievictor/ogugu && cd ogugu
```
2. Install dependencies
```bash
go get ./..
```
3. Copy the `.env.example` file to `.env` and provide the necessary values for your local setup:
```bash
cp .env.example .env
```
4. Set up the database and cache:
```bash
docker compose up -d
```
> This will expose PostgreSQL and Redis on your machine for Ogugu to connect to.
5. Run Database Migrations
```bash
export PGURL=<databaseurl>
migrate -database ${PGURL} -path migrations up
```

### Kubernetes and Monitoring Setup (optional)
1. Setup a kubernetes cluster
> For local development and testing, kind (Kubernetes IN Docker) is recommended. It provides a lightweight Kubernetes cluster running inside Docker containers on your machine.
```bash
kind create cluster
```
2. Provision infrastructure with Terraform
```bash
cd terraform
terraform init
terraform plan   
terraform apply
```
3. Verify Deployment
```bash
kubectl get pods -n monitoring
```
4. Expose monitoring dashboards to localhost
```bash
# forward grafana port
kubectl port-forward svc/grafana 3000:80 -n monitoring
# forward tempo port
kubectl port-forward svc/tempo 4318 -n monitoring
```

## Usage
After completing the setup and deployment, follow these steps to run the application and interact with its API:
1. Generate Swagger Docs
```bash
swag init
```
2. Run the application
```bash
go run main.go
```
3. Access the Swagger UI to test the API at `http://localhost:<port>/v1/swagger/index.html`
### View monitoring dashboards (if set up):
- Navigate to: http://localhost:3000
- Login credentials can be found in [terraform/values/grafana.yaml](./terraform/values/grafana.yaml) 
