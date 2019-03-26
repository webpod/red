workflow "Test" {
  on = "push"
  resolves = ["Golang Action"]
}

action "Golang Action" {
  uses = "cedrickring/golang-action@1.2.0"
}
