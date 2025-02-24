pipeline {
    agent any

    libraries {
        lib('fm-shared-library@main')
    }//end libraries. Github Repo: https://github.com/rollout/fm-cbci-shared-library

    options {
        // timestamps()
        timeout(time: 45, unit: 'MINUTES')
    }//end options

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage("Run unit tests"){
            agent {
                kubernetes {
                    label 'unit-test-' + UUID.randomUUID().toString()
                    inheritFrom 'default'
                    yamlFile './cbci-templates/fmforci.yaml'
                }
            }

            steps {
                container(name: "server", shell: "sh") {
                    withCredentials([
                        sshUserPrivateKey(credentialsId: 'SDK_E2E_SSH_KEY', keyFileVariable: 'SDK_E2E_SSH_KEY', passphraseVariable: '', usernameVariable: 'cloudbees.eslint@cloudbees.com'),
                        file(credentialsId: 'ENV_SECRETS', variable: 'ENV_SECRETS_PATH'),
                    ]) {
                        addGitHubFingerprint()
                        echo "====++++executing Run tests++++===="
                        sh script: 'cd ./v6 && /usr/local/go/bin/go test ./core/...', 
                           label: "Running unit tests"
                    }
                }
            }
            post{
                success{
                    script {
                        echo 'Unit Tests OK; posting results'
                        currentBuild.result = 'SUCCESS'
                    }
                }
                failure{
                    echo 'Unit Tests Failed;'
                }
            }
        }
    }
}
