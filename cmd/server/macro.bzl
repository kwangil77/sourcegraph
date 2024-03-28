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

        is_external = Label(target).repo_name != native.repo_name()
        if not is_external and len(kwargs) > 0:
            fail("non-external container dependency {target} shouldn't have additional kwargs, set them at the declaration site".format(target = target))

        target_lib = target + ":" + name + "_lib"
        target_bin = name + "_vbin" if is_external else name

        # If we're including a non workspace local dependency, we need to wrap it in our own
        # go_binary so we can apply attrs such as x_defs. For workspace local dependencies,
        # we should add it to the original go_binary and use that instead.
        if is_external:
            go_binary(
                name = target_bin,
                embed = [target_lib],
                visibility = ["//visibility:private"],
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
