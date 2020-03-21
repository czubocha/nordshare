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
	sam deploy -t build-output/backend.yaml --stack-name ${APP_NAME}-backend \
		--s3-bucket nordshare-pipeline-artifactstore-sacci84s97in --s3-prefix build-output  \
		--capabilities CAPABILITY_IAM