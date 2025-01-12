pipeline {
    agent any
    tools { go '1.23' }
    stages{
        stage("Cleanup"){
            steps {
                script {
                    def root = tool name: '1.23', type: 'go'
                    withEnv(["GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
                        sh '''
                            gofmt -s -w . && git diff --exit-code
                            go vet ./...
                            go mod tidy && git diff --exit-code
                            go mod download
                            go mod verify
                        '''
                    }
                }
                
            }
        }
        stage("Build"){
            steps {
                script {
                    def root = tool name: '1.23', type: 'go'
                    withEnv(["GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
                        sh "go build"
                    }
                }
            }
        }
    }
}