resource "google_pubsub_topic" "csp-demo-topic" {
  name = "csppubsubdemo-topic"
}

resource "google_pubsub_subscription" "csp-demo-subscription" {
  name  = "csppubsubdemo-subscription"
  topic = google_pubsub_topic.csp-demo-topic.id

  ack_deadline_seconds = 20

 labels = {
    foo = "bar"
  }

  push_config {
    push_endpoint = "https://csppubsubdemo-762434879017.us-central1.run.app/pubsub"

    attributes = {
      x-goog-version = "v1"
    }
  }
}

resource "google_pubsub_topic_iam_member" "publisher" {
  topic = google_pubsub_topic.csp-demo-topic.name
  role  = "roles/pubsub.publisher"
  member = "serviceAccount:cloud-support-apievents@prod.google.com"
}