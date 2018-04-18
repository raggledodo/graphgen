load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def dependencies():
    # tenncor dependency
    if "com_github_mingkaic_tenncor" not in native.existing_rules():
        git_repository(
            name = "com_github_mingkaic_tenncor",
            remote = "https://github.com/mingkaic/tenncor",
            commit = "9023ee19972dbab768195c2a2a7438da6b4f0476",
        )

    if "com_github_mingkaic_go_tenncor" not in native.existing_rules():
        git_repository(
            name = "com_github_mingkaic_go_tenncor",
            remote = "https://github.com/mingkaic/go_tenncor",
            commit = "d672a6b4020d861423172c01d540f62ab5d08af3",
        )

    # python dependency
    if "io_bazel_rules_python" not in native.existing_rules():
        git_repository(
            name = "io_bazel_rules_python",
            remote = "https://github.com/bazelbuild/rules_python.git",
            commit = "b25495c47eb7446729a2ed6b1643f573afa47d99",
        )

    # go dependency
    if "io_bazel_rules_go" not in native.existing_rules():
        http_archive(
            name = "io_bazel_rules_go",
            urls = [ "https://github.com/bazelbuild/rules_go/releases/download/0.10.3/rules_go-0.10.3.tar.gz" ],
            sha256 = "feba3278c13cde8d67e341a837f69a029f698d7a27ddbb2a202be7a10b22142a",
        )

    if "bazel_gazelle" not in native.existing_rules():
        http_archive(
            name = "bazel_gazelle",
            urls = [ "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.10.1/bazel-gazelle-0.10.1.tar.gz" ],
            sha256 = "d03625db67e9fb0905bbd206fa97e32ae9da894fe234a493e7517fd25faec914",
        )
