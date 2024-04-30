variable "sources" {
  type = list(map(string))
  default = [
    {
      filename = "test/foo.txt"
      contents = "foo"
    },
    {
      filename = "test/bar.txt"
      contents = "bar"
    }
  ]
}