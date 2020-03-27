.EXPORT_ALL_VARIABLES:

APP_NAME = nordshare

parent_deploy:
	sam validate -t deployments/parent.yml
	sam validate -t deployments/frontend.yml
	sam validate -t deployments/pipeline.yml
	sam deploy --template-file deployments/parent.yml --stack-name ${APP_NAME} \
		--s3-bucket ${APP_NAME}-templates \
		--capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND \
		--parameter-overrides ParameterKey=AppName,ParameterValue=${APP_NAME} \
        	ParameterKey=GitHubOwner,ParameterValue=czubocha \
           	ParameterKey=GitHubRepo,ParameterValue=${APP_NAME} \
          	ParameterKey=GitHubBranch,ParameterValue=master \
        	ParameterKey=SecretName,ParameterValue=${APP_NAME} \
        	ParameterKey=GitHubTokenSecretKey,ParameterValue=github-token \
        	ParameterKey=RefererSecretKey,ParameterValue=referer

backend_deploy:
	sam validate -t deployments/backend.yml
	sam build -t deployments/backend.yml -b build-output
	sam deploy -t build-output/template.yaml --stack-name ${APP_NAME}-backend \
		--s3-bucket nordshare-pipeline-r3jtocla6l2x-artifactstore-i7g9wmuxb77c --s3-prefix build-output  \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides ParameterKey=DeployerRoleArn,ParameterValue=arn:aws:iam::071572870590:role/nordshare-Pipeline-R3JTOCLA6L2X-DeployBackendRole-VMRWG4BSFMMP

cloudfront_deploy:
	sam validate -t deployments/cloudfront.yml
	sam deploy -t deployments/cloudfront.yml --stack-name ${APP_NAME}-cloudfront \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides ParameterKey=ApiId,ParameterValue=rhf5cvvuog \
			ParameterKey=ApiStage,ParameterValue=stage ParameterKey=RefererSecret,ParameterValue=$REFERER_SECRET \
			ParameterKey=FrontendWebsiteURL,ParameterValue=http://nordshare.s3-website.eu-central-1.amazonaws.com

frontend_deploy:
	cd website && npm run build && aws s3 sync build/ s3://nordshare

invoke_e2e_local:
	sam local invoke -t build-output/template.yaml -n configs/local-env.json E2E

start_api_local: build_handlers
	sam local start-api -t deployments/backend.yml -p 8080 -n configs/local-env.json

start_frontend_local:
	cd website && npm start

build_handlers:
	GOOS=linux go build -o cmd/saver/handler ./cmd/saver
	GOOS=linux go build -o cmd/reader/handler ./cmd/reader
	GOOS=linux go build -o cmd/modifier/handler ./cmd/modifier
	GOOS=linux go build -o cmd/remover/handler ./cmd/remover