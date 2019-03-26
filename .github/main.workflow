workflow "Build" {
  resolves = ["Golang Action"]
  on = "push"
}

action "Golang Action" {
  uses = "cedrickring/golang-action@1.2.0"
}
