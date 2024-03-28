load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

def get_last_segment(path):
    segments = path.split("/")
    last_segment = segments[-1]

    s = last_segment.split(":")
    if len(s) == 1:
        return last_segment
    else:
        return s[-1]

def container_dependencies(targets):
    for target, kwargs in targets.items():
        name = get_last_segment(target)

        target_lib = target + ":" + name + "_lib"
        target_bin = name + "_vbin"

        go_binary(
            name = target_bin,
            embed = [target_lib],
            visibility = ["//visibility:public"],
            **kwargs
        )

        pkg_tar(
            name = "tar_{}".format(name),
            srcs = [target_bin],
            remap_paths = {"/{}".format(name): "/usr/local/bin/{}".format(name)},
        )

def dependencies_tars(targets):
    tars = []
    for target, _ in targets.items():
        name = get_last_segment(target)
        tars.append(":tar_{}".format(name))

    return tars
