all:
	bazel run //server:server

update:
	glide update
	bazel run //:gazelle
