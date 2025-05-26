pipeline {
  agent any
  stages {
    stage('Hello') {
      steps {
        echo 'Hello from Jenkins!'
      }
    }
    stage('Check Kubernetes') {
      steps {
        sh 'kubectl get nodes'
      }
    }
  }
}
