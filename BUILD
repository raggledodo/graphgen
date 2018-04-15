licenses(["notice"])

package(
    default_visibility = ["//visibility:public"],
)

load("@pip_grpcio//:requirements.bzl", "requirement")
load("@bazel_gazelle//:def.bzl", "gazelle")

# graph generator
py_binary(
    name = "tfgen",
    srcs = ["tfgen.py"],
    deps = [
        requirement("grpcio"),
        "@com_github_mingkaic_tenncor//tests/py:tenncorgen",
        "@com_github_mingkaic_tenncor//tests/graphmgr:graphmgr_py_grpc",
    ],
)

# generator
gazelle(
    name = "gazelle",
    # external = "vendored",
    prefix = "github.com/raggledodo/graphmgr",
)
