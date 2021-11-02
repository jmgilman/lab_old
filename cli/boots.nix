{ lib, buildGoModule }:

buildGoModule rec {
    pname = "boots";
    version = "0.1.1";

    vendorSha256 = "1ffqyhc5ar51rsph7jahzzlf3qdgm5aq4y3fblwm5q7vqnriwbwy";

    src = ./.;
    subPackages = [ "./cmd/boots" ];
}