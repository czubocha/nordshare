pipeline_delete:
	aws cloudformation delete-stack --stack-name nordshare-pipeline

pipeline_deploy:
	sam validate -t deployments/pipeline.yml
	sam deploy --template-file deployments/pipeline.yml --stack-name nordshare-pipeline \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides ParameterKey=AppName,ParameterValue=nordshare \
	ParameterKey=GitHubOwner,ParameterValue=czubocha \
   	ParameterKey=GitHubRepo,ParameterValue=nordshare \
  	ParameterKey=GitHubBranch,ParameterValue=master \
	ParameterKey=GitHubOAuthTokenSecretName,ParameterValue=github \
	ParameterKey=GitHubOAuthTokenSecretKey,ParameterValue=token

backend_deploy:
	sam validate -t deployments/backend.yml
	sam build -t deployments/backend.yml -b build-output
	sam deploy -t build-output/backend.yaml --stack-name nordshare-backend \
		--s3-bucket nordshare-pipeline-artifactstore-sacci84s97in --s3-prefix build-output  \
		--capabilities CAPABILITY_IAM