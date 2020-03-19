pipeline delete:
	aws cloudformation delete-stack --stack-name nordshare
pipeline deploy:
	sam deploy --template-file deployments/pipeline.yml --stack-name nordshare \
	--capabilities CAPABILITY_IAM \
	--parameter-overrides ParameterKey=AppName,ParameterValue=nordshare \
	ParameterKey=GitHubOwner,ParameterValue=czubocha \
   	ParameterKey=GitHubRepo,ParameterValue=nordshare \
  	ParameterKey=GitHubBranch,ParameterValue=master \
	ParameterKey=GitHubOAuthTokenSecretName,ParameterValue=github \
	ParameterKey=GitHubOAuthTokenSecretKey,ParameterValue=token