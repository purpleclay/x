{
  lib,
  buildGoWorkspace,
  go,
}:
buildGoWorkspace {
  inherit go;

  pname = "theme";
  version = "0.1.0";
  src = ./.;
  subPackages = ["theme/cmd/theme"];
  modules = ./govendor.toml;
  ldflags = ["-s"];

  CGO_ENABLED = 0;

  meta = with lib; {
    description = "A Purple Clay theme used across all projects";
    homepage = "https://github.com/purpleclay/x";
    license = licenses.mit;
    maintainers = with maintainers; [purpleclay];
  };

  doCheck = false;
}
