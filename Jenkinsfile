@Library('main-shared-library') _

def git_repo_name = scm.getUserRemoteConfigs()[0].getUrl().replaceFirst(/^.*\/([^\/]+?).git$/, '$1')
def artifact_name = git_repo_name.replace('.', '-').toLowerCase()

properties([
	parameters([
		string(name: 'BASE_VERSION', defaultValue: '1.0')
	])
])

microservice_ci_go([
	github_repo_name: git_repo_name,
	base_image_uri: "534369319675.dkr.ecr.us-west-2.amazonaws.com/jnlp-test:latest",
	ecr_uri: "534369319675.dkr.ecr.us-west-2.amazonaws.com",
	artifact_name: artifact_name,
	dockerfile_context: "`pwd`/cmd/builder/output/slauth",
	dockerfile_path: "`pwd`/cmd/builder/output/slauth/Dockerfile",
	app_name: "sealights-risk-management-sl-collector",
	build_path: "cmd/builder",
    build_command: "./builder ./builder --config build-config.yaml --output-path=./output/slauth --name=sl-exporter",
	base_version: params.BASE_VERSION
])
