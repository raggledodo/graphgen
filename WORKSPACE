workspace(name = "com_github_mingkaic_graphgen")

load("//:graphmgr.bzl", "dependencies")
dependencies()

load("@com_github_mingkaic_tenncor//:tenncor.bzl", "dependencies", "test_dependencies")
dependencies()
test_dependencies()

# python dependencies
load("@io_bazel_rules_python//python:pip.bzl", "pip_repositories", "pip_import")
pip_repositories()
pip_import(
   name = "pip_grpcio",
   requirements = "@org_pubref_rules_protobuf//python:requirements.txt",
)

load("@pip_grpcio//:requirements.bzl", pip_grpcio_install = "pip_install")
pip_grpcio_install()

load("@org_pubref_rules_protobuf//python:rules.bzl", "py_proto_repositories")
py_proto_repositories()

# go dependencies
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

load("@org_pubref_rules_protobuf//go:rules.bzl", "go_proto_repositories")
go_proto_repositories()
