pipeline {
  agent any

  environment {
    REGISTRY = "docker.io/yogarn"
    USERNAME = "yogarn"
    IMAGE_NAME = "filkompedia-be"
    BUILD_NUMBER = "${env.BUILD_NUMBER}"
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
            echo "$PASS" | docker login --username="$USER" --password-stdin
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
	  kubectl apply -f env-configmap.yaml
          sed "s|__BUILD_NUMBER__|${BUILD_NUMBER}|g" ./kubernetes/deployment.yaml > ./kubernetes/deployment.generated.yaml
          kubectl apply -f ./kubernetes/deployment.generated.yaml
          kubectl apply -f ./kubernetes/service.yaml
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

