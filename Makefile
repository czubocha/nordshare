pipeline_delete:
	aws cloudformation delete-stack --stack-name nordshare-pipeline

pipeline_deploy:
	sam deploy --template-file deployments/pipeline.yml --stack-name nordshare-pipeline \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides ParameterKey=AppName,ParameterValue=nordshare \
	ParameterKey=GitHubOwner,ParameterValue=czubocha \
   	ParameterKey=GitHubRepo,ParameterValue=nordshare \
  	ParameterKey=GitHubBranch,ParameterValue=master \
	ParameterKey=GitHubOAuthTokenSecretName,ParameterValue=github \
	ParameterKey=GitHubOAuthTokenSecretKey,ParameterValue=token

deploy:
	sam build -t deployments/template.yml -b build-output
	sam deploy -t build-output/template.yaml --stack-name nordshare \
		--s3-bucket nordshare-pipeline-artifactstore-sacci84s97in --s3-prefix build-output  \
		--capabilities CAPABILITY_IAM