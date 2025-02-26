pipeline {
    agent any

    libraries {
        lib('fm-shared-library@main')
    }//end libraries. Github Repo: https://github.com/rollout/fm-cbci-shared-library

    options {
        timeout(time: 45, unit: 'MINUTES')
    }//end options

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage("Run Unit tests"){
            agent {
                kubernetes {
                    label 'unit-test-' + UUID.randomUUID().toString()
                    inheritFrom 'golang'
                    yamlFile './cbci-templates/fmforci.yaml'
                }
            }

            steps {
                container(name: "server", shell: "sh") {
                    withCredentials([
                        sshUserPrivateKey(credentialsId: 'SDK_E2E_SSH_KEY', keyFileVariable: 'SDK_E2E_SSH_KEY', passphraseVariable: '', usernameVariable: 'cloudbees.eslint@cloudbees.com'),
                    ]) {
                        withEnv(["PATH+GO=$PATH:/usr/local/go/bin"]) {
                            echo "Executing Run tests"
                            sh script: 'cd ./v6 && go test ./core/...', 
                                label: "Running unit tests"
                        }
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

        stage("Run E2E tests"){
            agent {
                kubernetes {
                    label 'e2e-tests-' + UUID.randomUUID().toString()
                    inheritFrom 'default'
                    yamlFile './cbci-templates/fmforci.yaml'
                }
            }

            when {
                branch 'master'
            }

            steps {
                container("rox-proxy") {
                    waitForRoxProxy()
                }

                container(name: "server", shell: 'sh') {
                    withCredentials([
                        string(credentialsId: 'TEST_E2E_BEARER', variable: 'TEST_E2E_BEARER'),
                        sshUserPrivateKey(credentialsId: 'SDK_E2E_SSH_KEY', keyFileVariable: 'SDK_E2E_SSH_KEY', passphraseVariable: '', usernameVariable: 'cloudbees.eslint@cloudbees.com'),
                        sshUserPrivateKey(credentialsId: 'SDK_E2E_TESTS_DEPLOY_KEY', keyFileVariable: 'SDK_E2E_TESTS_DEPLOY_KEY', passphraseVariable: '', usernameVariable: 'cloudbees.eslint@cloudbees.com'),                        
                    ]) {
                        script {
                            addGitHubFingerprint()
                            TESTENVPARAMS = "QA_E2E_BEARER=$TEST_E2E_BEARER API_HOST=https://api.test.rollout.io CD_API_ENDPOINT=https://api.test.rollout.io/device/get_configuration CD_S3_ENDPOINT=https://rox-conf.test.rollout.io/ SS_API_ENDPOINT=https://api.test.rollout.io/device/update_state_store/ SS_S3_ENDPOINT=https://rox-state.test.rollout.io/ CLIENT_DATA_CACHE_KEY=client_data ANALYTICS_ENDPOINT=https://analytic.test.rollout.io/ NOTIFICATIONS_ENDPOINT=https://push.test.rollout.io/sse"

                            withEnv(["GIT_SSH_COMMAND=ssh -i ${SDK_E2E_SSH_KEY}", "PATH+GO=$PATH:/usr/local/go/bin"]) {
                                echo "Executing E2E tests"
                                sh script: """
                                    apt-get update && apt-get install -y curl gnupg
                                    curl -sL https://deb.nodesource.com/setup_lts.x | bash - 
                                    apt-get install -y nodejs && npm install -g yarn
                                    
                                    git clone git@github.com:rollout/sdk-end-2-end-tests.git
                                    ls -la
                                    ln -s ./v6/driver ./sdk-end-2-end-tests/drivers/go
                                    cd sdk-end-2-end-tests
                                    yarn install --frozen-lockfile
                                    SDK_LANG=go ${TESTENVPARAMS} NODE_ENV=container yarn test:env
                                """, label: "Pull SDK end2 tests repository"
                            }// end withEnv
                        }
                    }
                }
            }
            post{
                success{
                    script {
                        echo 'E2E Tests OK; posting results'
                        currentBuild.result = 'SUCCESS'
                    }
                }
                failure{
                    echo 'E2E Tests Failed;'
                }
        
            }
        }
    }
}
