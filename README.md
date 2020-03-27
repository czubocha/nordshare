# :scroll: Nordshare
## Share notes securely. Simple. :sparkles:
### tech stack:
* [Go](https://golang.org) :bowtie: 
  * [AWS SDK for Go](https://github.com/aws/aws-sdk-go)
  * [envconfig](https://github.com/kelseyhightower/envconfig)
  * [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)
* [React](https://reactjs.org) :globe_with_meridians:
  * [React Router](https://reacttraining.com/react-router/web/guides/quick-start)
  * [Bootstrap](https://getbootstrap.com)
  * [React Bootstrap](https://react-bootstrap.github.io)
* [AWS](https://aws.amazon.com) :cloud:
  * Lambda
  * CloudFront
  * DynamoDB
  * API Gateway
  * S3
  * Key Management System
  * X-Ray
  * CodePipeline
  * CodeBuild
  * CodeDeploy
  * CloudFormation
  
### tools:
* [Docker](https://www.docker.com) :whale: 
* [AWS Serverless Application Model](https://github.com/awslabs/serverless-application-model) :squirrel:
* [Insomnia](https://insomnia.rest)
* [GNU Make](https://www.gnu.org/software/make) :ox:

### the fanciest features :tada:
<details>
<summary><b>ONE pipeline for backend and frontend with 0-click (master to production) continunous deployment with cacheable build environments, E2E tests and gradual deployment monitored by alarms with automatic rollback</b></summary>
  <ol>
    <li>Source code pulling - repository webook on master branch</li>
    <li>Concurrently building, executing tests and caching build environments for the future use</li>
    <li>Pre-traffic hook with E2E tests executing</li>
    <li>Deployment with gradual production traffic shifting (alarms definition for monitoring and automatic rollback possible)</li>
  </ol>
</details>
<details>
  <summary><b>FULL local developing & testing with frontend <--> backend communication</b></summary>
<ol>
  <li>Build all lambdas binaries for Linux environment at once with <code>sam build</code> command</li>
  <li>Start local API Gateway with endpoints handled by lambdas in Docker containers</li>
  <li>Start frontend serving on localhost with proxying request to API</li>
</ol>
</details>
<details>
  <summary><b>ONE-click pipeline, backend, frontend deployment from localhost with Makefile targets definitions thanks to ~1100 lines of CloudFormation YAML files (even automatic Github webhook creation) with REALLY least privilege IAM roles permissions (created through trial and error because of no other option)</b></summary>
</details>
<details>
  <summary><b>ONE definiton of API in OpenAPI 3.0 standard used for creation actual API Gateway infrastructure and at the same time for generating human-readable documentation (for example with <a href='https://github.com/Redocly/redoc'<a>ReDoc</a>)</b></summary>
</details>  

### interesting facts:
* AWS API Gateway lowercase headers for HTTP2 https://http2.github.io/http2-spec/#HttpHeaders'>https://http2.github.io/http2-spec/#HttpHeaders
* local AWS API Gateway (`sam local start-api`) capitalize headers

---
#### password rules:
##### read password:
* reading note & time to note expiration
##### write password:
* reading note & time to note expiration
* modifying content of note & time to note expiration
* deleting note

* _password setting (both read & write) is **not** required_
* _max TTL is 1 day (1440 minutes)_
* _extending TTL possible by modifying (except case when write password is not set)_
* _private key for encryption rotating yearly_

| read password set | write password set| read access                 | write access           
| :---:             |:---:              |:---:                        |:---:   
| √                 | √                 | w/ read or write password   | with write password
| √                 | X                 | w/ read password            | X
| X                 | √                 | w/o password (open)         | with write password
| X                 | X                 | w/o password (open)         | X
