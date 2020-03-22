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