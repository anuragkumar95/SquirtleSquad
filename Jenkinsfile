pipeline {
   agent any
   stages {
       stage('Install docker'){
           steps{
               sh '''
                      sudo apt --assume-yes install docker.io
                      sudo systemctl start docker
                      sudo systemctl enable docker
                      python -m pip uninstall -y urllib3
                      python -m pip install urllib3==1.22
                  '''
               }
        }

        stage('Build') {
            steps {
                sh 'echo "printing the build meter"'
                sh '''
                    echo "test"
                    ls -lah
                '''
                }
            }
        }

        stage('Building SessionManagement Service') {
            steps {
                dir('SessionManagement/') {
                       sh '''
                       docker build -t maxprimex123/squirtlesquad_sessionmanagement:latest .
                       docker push maxprimex123/squirtlesquad_sessionmanagement:latest
                       '''
                }
            }
        }

        stage('Building UserManagement Service') {
            steps {
                dir('/UserManagement/') {
                       sh '''
                       docker build -t maxprimex123/squirtlesquad_usermanagement:latest .
                       docker push maxprimex123/squirtlesquad_usermanagement:latest
                       '''
                }
            }
        }

        stage('Building APIGateway Service') {
            steps {
                dir('/APIGateway/') {
                       sh '''
                       docker build -t maxprimex123/squirtlesquad_apigateway:latest .
                       docker push maxprimex123/squirtlesquad_apigateway:latest
                       '''
                }
            }
        }

        stage('Building ModelExecution Service') {
            steps {
                dir('/ModelExecution/') {
                       sh '''
                       docker build -t maxprimex123/squirtlesquad_modelexecution:latest .
                       docker push maxprimex123/squirtlesquad_modelexecution:latest
                       '''
                }
            }
        }

        stage('Build PostAnalysis Service') {
            steps {
                dir('/Inference') {
                       sh '''
                       docker build -t maxprimex123/squirtlesquad_inferece:latest .
                       docker push maxprimex123/squirtlesquad_inference:latest
                       '''
                }
            }
        }


        stage('Deploy to K8s'){
            steps{
                dir('/Kubes2/') {
                    sh '''
                    sudo ssh  -i id_rsa ubuntu@149.165.171.111 &&
                    sudo apt install git -y &&
                    git clone https://github.com/airavata-courses/SquirtleSquad &&
                    cd SquirtleSquad &&
                    git checkout dockerized_services &&
                    cd SquirtleSquad/Kubes2/ &&
                    kubectl delete deployment,svc apigateway usermanagement sessionmanagement dataretrieval modelexecution inference &&
                    kubectl delete deployment zookeeper-dep kafka-dep &&
                    kubectl delete svc zookeeper kafka-service &&
                    kubectl apply -f message.yml &&
                    kubectl apply -f apigateway.yml &&
                    kubectl apply -f usermanagement.yml &&
                    kubectl apply -f sessionmanagement.yml &&
                    kubectl apply -f dataretrieval.yml &&
                    kubectl apply -f modelexecution.yml &&
                    kubectl apply -f inference.yml &&
                    kubectl scale deployment/apigateway deployment/usermanagement deployment/sessionmanagement deployment/dataretrieval deployment/modelexecution deployment/inference --replicas=3
                    '''
                }
            }
        }
    }
}
