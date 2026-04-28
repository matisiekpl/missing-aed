pipeline {
    agent any

    environment {
        IMAGE_NAME = "missing-aed"
        IMAGE_TAG = "latest"
        IMAGE_TAR = "missing-aed.tar"
        DEPLOYMENT_NAME = "missing-aed"
        K8S_NAMESPACE = "aed"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build image') {
            steps {
                sh 'docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .'
            }
        }

        stage('Export image') {
            steps {
                sh 'docker save ${IMAGE_NAME}:${IMAGE_TAG} -o ${IMAGE_TAR}'
            }
        }

        stage('Import to k3s containerd') {
            steps {
                sh 'k3s ctr images import ${IMAGE_TAR}'
            }
        }

        stage('Apply manifests') {
            steps {
                sh 'kubectl apply -f k8s.yaml'
            }
        }

        stage('Rollout') {
            steps {
                sh 'kubectl -n ${K8S_NAMESPACE} rollout restart deployment/${DEPLOYMENT_NAME}'
                sh 'kubectl -n ${K8S_NAMESPACE} rollout status deployment/${DEPLOYMENT_NAME} --timeout=180s'
            }
        }
    }

    post {
        always {
            sh 'rm -f ${IMAGE_TAR}'
        }
    }
}
