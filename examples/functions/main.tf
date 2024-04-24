locals {
  base64targz = (
    provider::functions::base64targz([
      {
        filename = "test/foo.txt"
        contents = "foo"
      },
      {
        filename = "test/bar.txt"
        contents = "bar"
      }
    ])
  )
}
