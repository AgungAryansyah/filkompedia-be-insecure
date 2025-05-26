pipeline {
  agent any

  environment {
    REGISTRY = "docker.io/yogarn"
    IMAGE_NAME = "filkompedia-be"
    BUILD_NUMBER = "${env.BUILD_NUMBER}"
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scm
      }
    }

    stage('Build Go App') {
      steps {
        sh '''
          go mod init filkompedia-be || true
          go mod tidy
          GOOS=linux GOARCH=amd64 go build -o main .
        '''
      }
    }

    stage('Build Docker Image') {
      steps {
        sh 'docker build -t ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} ./filkompedia-be'
      }
    }

    stage('Push Docker Image') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'docker-registry-creds', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
          sh '''
            echo "$PASS" | docker login ${REGISTRY} -u "$USER" --password-stdin
            docker push ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER}
            docker logout ${REGISTRY}
          '''
        }
      }
    }

    stage('Deploy to Kubernetes') {
      steps {
        sh '''
          sed "s|__BUILD_NUMBER__|${BUILD_NUMBER}|g" goapp/kubernetes/deployment.yaml > goapp/kubernetes/deployment.generated.yaml
          kubectl apply -f goapp/kubernetes/deployment.generated.yaml
          kubectl apply -f goapp/kubernetes/service.yaml
        '''
      }
    }
  }

  post {
    success {
      echo "✅ Deployed version ${BUILD_NUMBER} to Kubernetes!"
    }
    failure {
      echo "❌ Deployment failed."
    }
  }
}
