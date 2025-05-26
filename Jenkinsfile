pipeline {
  agent any

  environment {
    REGISTRY = "docker.io/yogarn"
    USERNAME = "yogarn"
    IMAGE_NAME = "filkompedia-be"
    BUILD_NUMBER = "yogarn${env.BUILD_NUMBER}"
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scm
      }
    }

    stage('Build Docker Image') {
      steps {
        sh 'docker build -t ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} .'
      }
    }

    stage('Push Docker Image') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'docker-registry-creds', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
          sh '''
            echo "$PASS" | docker login ${REGISTRY} -u "$USER" --password-stdin
	    docker tag ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER}
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

