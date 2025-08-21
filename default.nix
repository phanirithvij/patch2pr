{
  lib,
  buildGoModule,
  testers,
  gitMinimal,
  makeWrapper,
  patch2pr,
}:

buildGoModule (finalAttrs: {
  pname = "patch2pr";
  version = "0.38.0";

  src = lib.cleanSource ./.;
  vendorHash = "sha256-QEgGq5/JQUIWWmJKoQ832eKhiF5xF8Jivpn1uFDERTA=";

  ldflags = [
    "-X main.version=${finalAttrs.version}"
    "-X main.commit=dev"
  ];

  nativeBuildInputs = [ makeWrapper ];

  postInstall = ''
    wrapProgram $out/bin/patch2pr --prefix PATH : ${lib.makeBinPath [ gitMinimal ]}
  '';

  passthru.tests.patch2pr-version = testers.testVersion {
    package = patch2pr;
    command = "${patch2pr.meta.mainProgram} --version";
    version = finalAttrs.version;
  };

  meta = {
    changelog = "https://github.com/bluekeyes/patch2pr/releases/tag/v${finalAttrs.version}";
    description = "Create pull requests from patches without cloning the repository";
    homepage = "https://github.com/bluekeyes/patch2pr";
    license = lib.licenses.mit;
    mainProgram = "patch2pr";
    maintainers = with lib.maintainers; [ phanirithvij ];
  };
})
